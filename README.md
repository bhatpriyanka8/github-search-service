# GitHub Search Service

A simple gRPC service built in Go on top of the GitHub Search API to search based on search phrase given, with optional filtering by user.

## Overview
This service exposes a gRPC API that accepts a search term and an optional GitHub username, performs a search using the GitHub REST Search API internally, and returns the file URLs and repositories where matches were found.

## API
```proto
service GithubSearchService { 
    rpc Search(SearchRequest) returns (SearchResponse); 
}