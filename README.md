# GitHub Search Service

A simple gRPC service written in Go that wraps the GitHub REST Code Search API.
It allows searching for code across GitHub repositories with optional filtering by user.


## What this service does

- Exposes a gRPC API to search GitHub code

- Accepts a search term and optional GitHub username

- Calls GitHub’s REST Search API internally

- Returns file URLs and repository names

- Uses proper gRPC status codes and input validation

## API (gRPC)

```
service GithubSearchService {
  rpc Search(SearchRequest) returns (SearchResponse);
}

message SearchRequest {
  string search_term = 1;
  string user = 2;
}

message SearchResponse {
  repeated Result results = 1;
}

message Result {
  string file_url = 1;
  string repo = 2;
}

```

## How to run

### Prerequisites

- Go 1.21+
- GitHub Personal Access TOken

### Authentication

This service uses the GitHub Code Search API, which requires authentication for reliable access.
Set a GitHub Personal Access Token:

```bash
export GITHUB_TOKEN=<your_token>
```

### Run the service

```go run ./cmd/service/main.go```

The server listens on localhost:8182.

**Example request (using grpcurl)**

```bash
grpcurl -plaintext \
  -d '{"search_term":"golang"}' \
  localhost:8182 \
  githubsearch.GithubSearchService/Search
  ```

**With user filter:**

```bash
grpcurl -plaintext \
  -d '{"search_term":"golang","user":"priyanka"}' \
  localhost:8182 \
  githubsearch.GithubSearchService/Search
  ```

## Project structure (high level)
cmd/service        # gRPC server entrypoint (main)
internal/service   # gRPC service implementation
internal/github    # GitHub REST API client
internal/gen       # Generated protobuf code

## Testing

Run all unit tests:

go test ./...

**Tests focus on:**

- Input validation
- Error handling
- Service logic using mocked GitHub client
- No real GitHub API calls are made in tests.

## Design notes & trade-offs

- Pagination is not implemented yet; the service returns the first page of results.
- GitHub rate limiting is respected. Respective Error is shown.
- Context timeouts are applied to downstream GitHub API calls.
- Authentication is supported via environment variable (GITHUB_TOKEN).

These choices were made to keep the implementation simple and focused.
