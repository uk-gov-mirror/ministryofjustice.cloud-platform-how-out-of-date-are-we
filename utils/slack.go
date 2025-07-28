package utils

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var (
	slackToken = os.Getenv("SLACK_BOT_TOKEN")
	channelID  = os.Getenv("SLACK_CHANNEL_ID")
)

func slackClient(token string) (*slack.Client, error) {
	if token == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN environment variable is not set")
	}

	client := slack.New(token)
	return client, nil
}

func postToChannel(channelName, message string) error {
	client, err := slackClient(os.Getenv("SLACK_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("error creating Slack client: %w", err)
	}

	_, _, err = client.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		return fmt.Errorf("error posting message to channel %s: %w", channelName, err)
	}

	return nil
}
