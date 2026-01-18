package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	token      string
}

// SearchResult is the result returned by the github client
type SearchResult struct {
	File string
	Repo string
}

// NewClient creates a Github Client instance
func NewClient(httpClient *http.Client, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		token:      token,
	}
}

// Search searches GitHub code using the provided search term and user if given
func (c *Client) Search(
	ctx context.Context,
	term string,
	user string,
) ([]SearchResult, error) {

	var builder strings.Builder
	var url = "https://api.github.com/search/code"
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(" github api returned status code %d", resp.StatusCode)
	}

	var githubResp githubSearchResponse

	err = json.NewDecoder(resp.Body).Decode(&githubResp)
	if err != nil {
		return nil, err
	}

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
