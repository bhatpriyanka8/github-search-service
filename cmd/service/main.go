package main

import (
	"fmt"
	"log"
	"net"
	"os"

	servicePB "github.com/bhatpriyanka8/github-search-service/internal/gen/githubsearch"
	"github.com/bhatpriyanka8/github-search-service/internal/github"
	"github.com/bhatpriyanka8/github-search-service/internal/service"
	"google.golang.org/grpc"
)

var Port = "8182"

func main() {

	log.Println("Starting GitHub Search Service...")

	// listen on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	grpcServer := grpc.NewServer()

	// create service
	token := os.Getenv("GITHUB_TOKEN")
	baseURL := os.Getenv("GITHUB_API_URL")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	githubClient := github.NewClient(nil, token, baseURL)
	searchService := service.NewGitHubSearchService(githubClient)

	// register service with grpc server
	servicePB.RegisterGithubSearchServiceServer(grpcServer, searchService)

	// start
	log.Printf("grpc server for github search listening on: %s", Port)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
