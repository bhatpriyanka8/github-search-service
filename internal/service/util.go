package service

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapAndConvertError(err error) error {
	errMsg := err.Error()

	// authentication errors (401)
	if strings.Contains(errMsg, "status code 401") {
		return status.Error(codes.PermissionDenied, "authentication failed: please set GITHUB_TOKEN environment variable with a valid personal access token")
	}

	// rate limit errors (403)
	if strings.Contains(errMsg, "status code 403") {
		return status.Error(codes.ResourceExhausted, "github api rate limit exceeded, please try again after a minute")
	}

	// internal error
	return status.Error(codes.Internal, fmt.Sprintf("github search failed: %s", errMsg))
}
