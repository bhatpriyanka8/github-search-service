package github

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Retry retries the input function for the specified attempts and time period, with exponential backoff
func Retry(ctx context.Context, maxTries int, waitTime time.Duration, f func() error) error {
	for i := 0; i < maxTries; i++ {
		err := f()
		if err == nil {
			return nil
		}
		// fail fast on rate limit
		if err.Error() == "rate limit exceeded" {
			return err
		}
		log.Printf("attempt %d/%d failed: %v", i+1, maxTries, err)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		time.Sleep(waitTime)
		waitTime *= 2 // Exponential backoff
	}
	return fmt.Errorf("failed after %d attempts", maxTries)
}
