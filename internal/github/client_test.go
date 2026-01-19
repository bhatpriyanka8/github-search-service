package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
		token      string
		baseURL    string
	}{
		{
			name:       "create client with http client and token",
			httpClient: &http.Client{},
			token:      "test-token",
			baseURL:    "",
		},
		{
			name:       "create client with nil http client",
			httpClient: nil,
			token:      "test-token",
			baseURL:    "api/github.com",
		},
		{
			name:       "create client without token",
			httpClient: &http.Client{},
			token:      "",
			baseURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.httpClient, tt.token, tt.baseURL)
			if client == nil {
				t.Error("NewClient() returned nil")
			}
		})
	}
}
func TestClient_Search(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		err        error
		term       string
		user       string
		ctxCancel  bool
		wantErr    bool
		wantFiles  []string
		wantRepos  []string
	}{
		{
			name:       "success",
			statusCode: 200,
			body:       `{"items":[{"html_url":"url1","repository":{"full_name":"repo1"}},{"html_url":"url2","repository":{"full_name":"repo2"}}]}`,
			term:       "grpc",
			user:       "testuser",
			wantFiles:  []string{"url1", "url2"},
			wantRepos:  []string{"repo1", "repo2"},
		},
		{
			name:       "empty results",
			statusCode: 200,
			body:       `{"items":[]}`,
			term:       "grpc",
			user:       "testuser",
			wantFiles:  []string{},
			wantRepos:  []string{},
		},
		{
			name:       "server error triggers retry",
			statusCode: 500,
			body:       `internal error`,
			term:       "grpc",
			user:       "testuser",
			wantErr:    true,
		},
		{
			name:       "bad request",
			statusCode: 422,
			body:       `bad request`,
			term:       "",
			user:       "",
			wantErr:    true,
		},
		{
			name:       "context cancelled",
			statusCode: 200,
			body:       `{"items":[]}`,
			term:       "grpc",
			user:       "testuser",
			ctxCancel:  true,
			wantErr:    true,
		},
		{
			name:       "rate limit exceeded",
			statusCode: 403,
			body:       `rate limit exceeded`,
			term:       "grpc",
			user:       "testuser",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			mockHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer mockHttp.Close()
			httpClient := mockHttp.Client()
			client := NewClient(httpClient, "", mockHttp.URL)
			ctx := context.Background()
			if tt.ctxCancel {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			results, err := client.Search(ctx, tt.term, tt.user)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(results) != len(tt.wantFiles) {
				t.Errorf("expected %d results, got %d", len(tt.wantFiles), len(results))
			}
			for i, res := range results {
				if res.File != tt.wantFiles[i] || res.Repo != tt.wantRepos[i] {
					t.Errorf("unexpected result at %d: got %+v, want file=%s repo=%s", i, res, tt.wantFiles[i], tt.wantRepos[i])
				}
			}
		})
	}
}
