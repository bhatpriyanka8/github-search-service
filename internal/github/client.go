//go:generate mockgen -source ./client.go -package github -destination ./client_mock.go
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GitHubSearcher interface
type GitHubSearcher interface {
	Search(ctx context.Context, term, user string) ([]SearchResult, error)
}

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// SearchResult is the result returned by the github client
type SearchResult struct {
	File string
	Repo string
}

// NewClient creates a Github Client instance
func NewClient(httpClient *http.Client, token string, baseURL string) GitHubSearcher {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		token:      token,
		baseURL:    baseURL,
	}
}

// Search searches GitHub code using the provided search term and user if given
// Returns an empty response if no results are found
func (c *Client) Search(
	ctx context.Context,
	term string,
	user string,
) ([]SearchResult, error) {

	var builder strings.Builder
	url := c.baseURL + "/search/code"
	headers := map[string]string{
		"Accept":     "application/vnd.github.v3+json",
		"User-Agent": "github-search-service",
	}

	if c.token != "" {
		headers["Authorization"] = "token " + c.token
	}

	builder.WriteString(term)

	if user != "" {
		builder.WriteString(" user:")
		builder.WriteString(user)
	}

	query := builder.String()

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	q := req.URL.Query()
	q.Set("q", query)
	req.URL.RawQuery = q.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	var err error
	retryErr := Retry(ctx, 3, 100*time.Millisecond, func() error {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}
		// Fail fast on rate limit (403)
		if resp.StatusCode == http.StatusForbidden {
			resp.Body.Close()
			return fmt.Errorf("rate limit exceeded")
		}
		return nil
	})
	if retryErr != nil {
		return nil, retryErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status code %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var githubResp githubSearchResponse

	err = json.NewDecoder(resp.Body).Decode(&githubResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode github response: %w", err)
	}

	// if no items found, return empty
	if len(githubResp.Items) == 0 {
		return []SearchResult{}, nil
	}

	// prepare the results response
	results := make([]SearchResult, 0, len(githubResp.Items))

	for _, item := range githubResp.Items {
		results = append(results, SearchResult{
			File: item.HTMLURL,
			Repo: item.Repository.FullName,
		})
	}

	return results, nil

}

type githubSearchResponse struct {
	Items []githubSearchItem `json:"items"`
}

type githubSearchItem struct {
	HTMLURL    string           `json:"html_url"`
	Repository githubRepository `json:"repository"`
}

type githubRepository struct {
	FullName string `json:"full_name"`
}
