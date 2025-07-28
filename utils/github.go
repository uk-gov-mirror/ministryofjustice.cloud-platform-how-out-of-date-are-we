package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v72/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

func AppClient(key, appid, installid string) (*github.Client, error) {
	privateKey := []byte(key)

	appIDInt, err := strconv.ParseInt(appid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting app ID to int64: %w", err)
	}

	installIDInt, err := strconv.ParseInt(installid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting installation ID to int64: %w", err)
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(appIDInt, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating application token source: %w", err)
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installIDInt, appTokenSource)

	oauthHttpClient := oauth2.NewClient(context.Background(), installationTokenSource)

	client := github.NewClient(oauthHttpClient)
	return client, nil
}

func CreateIssueWithUpdateDetails(key, appid, installid, org, repo string, issueDetails []string) ([]string, int, error) {
	if key == "" || appid == "" || installid == "" {
		return nil, 0, fmt.Errorf("GITHUB_APP_KEY, GITHUB_APP_ID, or GITHUB_INSTALLATION_ID environment variables are not set")
	}

	githubClient, err := AppClient(key, appid, installid)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating GitHub client: %w", err)
	}
	ctx := context.Background()

	issue, _, err := githubClient.Issues.Create(ctx, org, repo, &github.IssueRequest{
		Title:  github.Ptr("Alert Manager Update"),
		Body:   github.Ptr(fmt.Sprint(issueDetails)),
		Labels: &[]string{"support", "slack receiver"},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error creating issue: %w", err)
	}

	issueOutput := []string{
		"Success: Issue created successfully",
		"IssueURL: " + issue.GetHTMLURL(),
		"IssueNumber: " + strconv.Itoa(issue.GetNumber()),
	}

	return issueOutput, issue.GetNumber(), nil
}

func CreatePullRequestForIssue(key, appid, installid, org, repo string, issue []string) ([]string, error) {
	if key == "" || appid == "" || installid == "" {
		return nil, fmt.Errorf("GITHUB_APP_KEY, GITHUB_APP_ID, or GITHUB_INSTALLATION_ID environment variables are not set")
	}

	githubClient, err := AppClient(key, appid, installid)
	if err != nil {
		return nil, fmt.Errorf("error creating GitHub client: %w", err)
	}
	ctx := context.Background()

	pr, _, err := githubClient.PullRequests.Create(ctx, org, repo, &github.NewPullRequest{
		Title: github.Ptr("Update Alert Manager"),
		Head:  github.Ptr("feature/alert-manager-update/" + issue[2]),
		Base:  github.Ptr("main"),
		Body:  github.Ptr("This pull request is created to address the issue related to Alert Manager updates. Please review the changes." + "\n\nIssue URL: " + issue[1]),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating pull request: %w", err)
	}

	prOutput := []string{
		"Success: Pull request created successfully",
		"PullRequestURL: " + pr.GetHTMLURL(),
	}

	return prOutput, nil
}
