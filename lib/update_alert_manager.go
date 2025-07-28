package lib

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/cloud-platform-how-out-of-date-are-we/utils"
)

var (
	s3Key = "alert_manager_receivers.json"
	// pKey      = os.Getenv("GITHUB_APP_KEY")
	// appid     = os.Getenv("GITHUB_APP_ID")
	// installid = os.Getenv("GITHUB_INSTALLATION_ID")
	// org       = "ministryofjustice"
	// repo      = "cloud-platform"
)

type AlertManagerUpdateResponseMulti []struct {
	Success  int
	Severity string
}

func UpdateAlertManager(w http.ResponseWriter, r *http.Request, s3bucket string, s3client *s3.Client, newAlert utils.AlertManagerUpdate, t *template.Template) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	file, err := utils.GetS3FileContent(s3bucket, s3Key, s3client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get S3 file content: %v", err), http.StatusInternalServerError)
		return
	}

	s3Map, err := utils.ParseAlertManagerReceiversFile(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse alert manager receivers file: %v", err), http.StatusInternalServerError)
		return
	}

	updatedReceivers, compareResult := utils.Compare(s3Map, newAlert)
	for _, res := range compareResult {
		if res < 0 || res > 2 {
			http.Error(w, "Invalid comparison result", http.StatusInternalServerError)
			return
		}
	}

	lines, err := utils.UpdateAlertDetails(updatedReceivers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update alert details: %v", err), http.StatusInternalServerError)
		return
	}

	if err := utils.UpdateS3FileDetails(s3bucket, s3Key, lines, s3client); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update S3 file: %v", err), http.StatusInternalServerError)
		return
	}

	var response AlertManagerUpdateResponseMulti
	response = AlertManagerUpdateResponseMulti{}
	for i, alert := range newAlert.Alert {
		if i < len(compareResult) { // Safety check
			response = append(response, struct {
				Success  int
				Severity string
			}{
				Success:  compareResult[i],
				Severity: alert.Severity,
			})
		}
	}

	templateData := struct {
		IssueURL                        string
		PullRequestURL                  string
		AlertManagerUpdateResponseMulti AlertManagerUpdateResponseMulti
	}{
		IssueURL:                        "feature coming soon",
		PullRequestURL:                  "feature coming soon",
		AlertManagerUpdateResponseMulti: response,
	}

	if err := t.Execute(w, templateData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}
