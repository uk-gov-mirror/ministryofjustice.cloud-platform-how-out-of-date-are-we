package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AlertManagerUpdate struct {
	TeamName string
	Contact  string
	Alert    []struct {
		Severity        string
		SlackWebhookURL string
		ChannelName     string
	}
}

type AlertManagerReceivers []struct {
	Severity        string `json:"severity"`
	SlackWebhookURL string `json:"webhook"`
	ChannelName     string `json:"channel"`
}

func CollectAlertManagerUpdates(r *http.Request, w http.ResponseWriter) (AlertManagerUpdate, error) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return AlertManagerUpdate{}, fmt.Errorf("unable to parse form: %w", err)
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return AlertManagerUpdate{}, fmt.Errorf("unsupported content type: %s", r.Header.Get("Content-Type"))
	}

	alerts := []struct {
		Severity        string
		SlackWebhookURL string
		ChannelName     string
	}{}

	channelNames := r.Form["ChannelName"]
	severities := r.Form["Severity"]
	webhookURLs := r.Form["SlackWebhookURL"]

	for i := range channelNames {
		alert := struct {
			Severity        string
			SlackWebhookURL string
			ChannelName     string
		}{
			Severity:        channelNames[i] + "-" + severities[i],
			SlackWebhookURL: webhookURLs[i],
			ChannelName:     "#" + channelNames[i],
		}
		alerts = append(alerts, alert)
	}

	newAlert := AlertManagerUpdate{
		TeamName: r.FormValue("TeamName"),
		Contact:  r.FormValue("Contact"),
		Alert:    alerts,
	}

	return newAlert, nil
}

func ParseAlertManagerReceiversFile(fileContent []byte) (AlertManagerReceivers, error) {
	var s3Content AlertManagerReceivers

	err := json.Unmarshal(fileContent, &s3Content)
	if err != nil {
		return AlertManagerReceivers{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	if len(s3Content) == 0 {
		return AlertManagerReceivers{}, fmt.Errorf("no alert manager receivers found in the file")
	}
	return s3Content, nil
}

func Compare(s3Content AlertManagerReceivers, newAlert AlertManagerUpdate) (AlertManagerReceivers, []int) {
	var compareResult []int
	for _, alert := range newAlert.Alert {
		found := false
		for newAlertIndex, s3Alert := range s3Content {
			if s3Alert.ChannelName == alert.ChannelName && s3Alert.Severity == alert.Severity {
				found = true
				if s3Alert.SlackWebhookURL == alert.SlackWebhookURL {
					compareResult = append(compareResult, 2)
				} else {
					s3Content[newAlertIndex].SlackWebhookURL = alert.SlackWebhookURL
					compareResult = append(compareResult, 1)
				}
				break
			}
		}
		if !found {
			fmt.Printf("Adding new alert: %s, %s, %s\n", alert.ChannelName, alert.SlackWebhookURL, alert.Severity)
			newReceiver := struct {
				Severity        string `json:"severity"`
				SlackWebhookURL string `json:"webhook"`
				ChannelName     string `json:"channel"`
			}{
				ChannelName:     alert.ChannelName,
				SlackWebhookURL: alert.SlackWebhookURL,
				Severity:        alert.Severity,
			}
			s3Content = append(s3Content, newReceiver)
			compareResult = append(compareResult, 0)
		}
	}

	return s3Content, compareResult
}

func UpdateAlertDetails(s3Map AlertManagerReceivers) ([]string, error) {
	if len(s3Map) == 0 {
		return []string{"[", "]"}, nil
	}

	lines := []string{}
	lines = append(lines, "[")

	for i, receiver := range s3Map {
		lines = append(lines, "  {")
		lines = append(lines, fmt.Sprintf(`    "severity": "%s",`, receiver.Severity))
		lines = append(lines, fmt.Sprintf(`    "webhook": "%s",`, receiver.SlackWebhookURL))
		lines = append(lines, fmt.Sprintf(`    "channel": "%s"`, receiver.ChannelName))

		if i == len(s3Map)-1 {
			lines = append(lines, "  }")
		} else {
			lines = append(lines, "  },")
		}
	}

	lines = append(lines, "]")
	return lines, nil
}

func GetJSONFileContent(s3Key string) ([]byte, error) {
	file, err := os.ReadFile("lib/static/development/" + s3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to read local file: %w", err)
	}
	return file, nil
}

func UpdateJSONFile(s3Key string, lines []string) error {
	// For testing - write to local file
	err := os.WriteFile("lib/static/development/"+s3Key, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated content to file: %w", err)
	}
	return nil
}
