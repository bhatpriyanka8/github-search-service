package service

import (
	"context"

	gs "github.com/bhatpriyanka8/github-search-service/internal/gen/githubsearch"
	github "github.com/bhatpriyanka8/github-search-service/internal/github"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GithubSearchService implements grpc service
type GitHubSearchService struct {
	githubClient *github.Client
	gs.UnimplementedGithubSearchServiceServer
}

// NewGitHubSearchService creates a new GitHubSearchService instance
func NewGitHubSearchService(ghClient *github.Client) *GitHubSearchService {
	return &GitHubSearchService{
		githubClient: ghClient,
	}
}

// Search implementation for doing github files/repos search
func (s *GitHubSearchService) Search(ctx context.Context, req *gs.SearchRequest) (*gs.SearchResponse, error) {
	if req.GetSearchTerm() == "" {
		return nil, status.Error(codes.InvalidArgument, "search_term missing, please provide")
	}
	results, err := s.githubClient.Search(ctx, req.GetSearchTerm(), req.GetUser())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "github search call failed: %v", err)
	}

	resp := &gs.SearchResponse{}

	for _, r := range results {
		resp.Results = append(resp.Results, &gs.Result{
			FileUrl: r.File,
			Repo:    r.Repo,
		})
	}

	return resp, nil
}
