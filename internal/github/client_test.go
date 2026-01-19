package github

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
		token      string
	}{
		{
			name:       "create client with http client and token",
			httpClient: &http.Client{},
			token:      "test-token",
		},
		{
			name:       "create client with nil http client",
			httpClient: nil,
			token:      "test-token",
		},
		{
			name:       "create client without token",
			httpClient: &http.Client{},
			token:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.httpClient, tt.token)
			if client == nil {
				t.Error("NewClient() returned nil")
			}
		})
	}
}
