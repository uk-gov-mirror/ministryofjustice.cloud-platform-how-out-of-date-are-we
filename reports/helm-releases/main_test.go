package main

import (
	"reflect"
	"testing"

	"k8s.io/client-go/kubernetes"
)

func Test_getCPNamespaces(t *testing.T) {
	type args struct {
		clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Get namespaces with specific annotation",
			args: args{
				clientset: &kubernetes.Clientset{}, // Mock or use a real clientset as needed
			},
			want:    []string{"namespace1", "namespace2"}, // Expected namespaces with the annotation
			wantErr: false,                                // Adjust based on your test setup
		},
		{
			name: "Error case",
			args: args{
				clientset: nil, // Simulate an error by passing nil or an invalid clientset
			},
			want:    nil,  // Expecting no namespaces
			wantErr: true, // Expecting an error due to nil clientset
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCPNamespaces(tt.args.clientset)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCPNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCPNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHelmReleasesInNamespaces(t *testing.T) {
	type args struct {
		namespaces []string
	}
	tests := []struct {
		name    string
		args    args
		want    []helmRelease
		wantErr bool
	}{
		{
			name: "Get helm releases in namespaces",
			args: args{
				namespaces: []string{"namespace1", "namespace2"},
			},
			want: []helmRelease{
				{
					Name:      "release1",
					Namespace: "namespace1",
					Chart:     "chart1",
				},
			},
			wantErr: false, // Adjust based on your test setup
		},
		{
			name: "Error case",
			args: args{
				namespaces: []string{"invalid-namespace"},
			},
			want:    nil,  // Expecting no releases
			wantErr: true, // Expecting an error due to invalid namespace
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHelmReleasesInNamespaces(tt.args.namespaces)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHelmReleasesInNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHelmReleasesInNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
