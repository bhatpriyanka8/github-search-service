package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gs "github.com/bhatpriyanka8/github-search-service/internal/gen/githubsearch"
	github "github.com/bhatpriyanka8/github-search-service/internal/github"
	"github.com/bhatpriyanka8/github-search-service/internal/service"
	"github.com/golang/mock/gomock"
)

func TestGitHubSearchService_Search(t *testing.T) {

	tests := []struct {
		name     string
		mockPrep func(mockClient *github.MockGitHubSearcher)
		req      *gs.SearchRequest
		want     *gs.SearchResponse
		wantErr  bool
	}{
		{
			name: "success",
			mockPrep: func(mockClient *github.MockGitHubSearcher) {
				mockClient.EXPECT().Search(gomock.Any(), "file", "").Return([]github.SearchResult{
					{
						File: "github.com/xyz/file.md",
						Repo: "xyz",
					},
					{
						File: "github.com/abc/file2.md",
						Repo: "abc",
					},
				}, nil)
			},
			req: &gs.SearchRequest{SearchTerm: "file"},
			want: &gs.SearchResponse{
				Results: []*gs.Result{
					{
						FileUrl: "github.com/xyz/file.md",
						Repo:    "xyz",
					},
					{
						FileUrl: "github.com/abc/file2.md",
						Repo:    "abc",
					},
				},
			},
		},
		{
			name:     "missing search term",
			mockPrep: func(mockClient *github.MockGitHubSearcher) {},
			req:      &gs.SearchRequest{SearchTerm: ""},
			wantErr:  true,
		},
		{
			name: "github client call failed",
			mockPrep: func(mockClient *github.MockGitHubSearcher) {
				mockClient.EXPECT().Search(gomock.Any(), "file", "").Return(nil, errors.New("internal error"))
			},
			req:     &gs.SearchRequest{SearchTerm: "file"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGHClient := github.NewMockGitHubSearcher(ctrl)
			tt.mockPrep(mockGHClient)
			s := service.NewGitHubSearchService(mockGHClient)

			got, gotErr := s.Search(context.Background(), tt.req)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Search() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Search() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
