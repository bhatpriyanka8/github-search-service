package service

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapAndConvertError(err error) error {
	if strings.Contains(err.Error(), "status code 401") {
		return status.Error(codes.PermissionDenied, "unauthorized GitHub API request, export valid github token")
	}
	if strings.Contains(err.Error(), "status code 403") {
		return status.Error(codes.ResourceExhausted, "GitHub rate limit exceeded, try after a minute")
	}
	return status.Error(codes.Internal, "GitHub search failed due to internal error")
}
