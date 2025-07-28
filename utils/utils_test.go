package utils

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCollectAlertManagerUpdates(t *testing.T) {
	type args struct {
		r *http.Request
		w http.ResponseWriter
	}
	tests := []struct {
		name    string
		args    args
		want    AlertManagerUpdate
		wantErr bool
	}{
		{
			name: "valid request",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/x-www-form-urlencoded"},
					},
					Form: map[string][]string{
						"ChannelName":     {"general"},
						"Severity":        {"critical"},
						"SlackWebhookURL": {"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"},
					},
				},
				w: http.ResponseWriter(http.ResponseWriter(nil)), // Mock or use a real ResponseWriter
			},
			want: AlertManagerUpdate{
				TeamName: "default",
				Contact:  "<contact>",
				Alert: []struct {
					Severity        string
					SlackWebhookURL string
					ChannelName     string
				}{
					{
						Severity:        "critical",
						SlackWebhookURL: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
						ChannelName:     "general",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid content type",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Form: map[string][]string{
						"ChannelName":     {"general"},
						"Severity":        {"critical"},
						"SlackWebhookURL": {"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"},
					},
				},
				w: http.ResponseWriter(http.ResponseWriter(nil)), // Mock or use a real ResponseWriter
			},
			want:    AlertManagerUpdate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CollectAlertManagerUpdates(tt.args.r, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectAlertManagerUpdates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAlertManagerUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseAlertManagerReceiversFile(t *testing.T) {
	type args struct {
		fileContent []byte
	}
	tests := []struct {
		name    string
		args    args
		want    AlertManagerReceivers
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				fileContent: []byte(`[{"severity": "critical", "webhook": "https://example.com/webhook", "channel": "alerts"}]`),
			},
			want: AlertManagerReceivers{
				{
					Severity:        "critical",
					SlackWebhookURL: "https://example.com/webhook",
					ChannelName:     "alerts",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				fileContent: []byte(`{"severity": "critical", "webhook": "https://example.com/webhook", "channel": "alerts"`),
			},
			want:    AlertManagerReceivers{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAlertManagerReceiversFile(tt.args.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAlertManagerReceiversFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAlertManagerReceiversFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		s3Content AlertManagerReceivers
		newAlert  AlertManagerUpdate
	}
	tests := []struct {
		name  string
		args  args
		want  AlertManagerReceivers
		want1 []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Compare(tt.args.s3Content, tt.args.newAlert)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compare() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Compare() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_UpdateAlertDetails(t *testing.T) {
	type args struct {
		s3Map AlertManagerReceivers
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid alert manager receivers",
			args: args{
				s3Map: AlertManagerReceivers{
					{
						Severity:        "critical",
						SlackWebhookURL: "https://example.com/webhook",
						ChannelName:     "alerts",
					},
				},
			},
			want: []string{
				"[",
				"  {",
				"    \"severity\": \"critical\",",
				"    \"webhook\": \"https://example.com/webhook\",",
				"    \"channel\": \"alerts\"",
				"  }",
				"]",
			},
			wantErr: false,
		},
		{
			name: "empty alert manager receivers",
			args: args{
				s3Map: AlertManagerReceivers{},
			},
			want:    []string{"[", "]"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateAlertDetails(tt.args.s3Map)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateAlertDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateAlertDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
