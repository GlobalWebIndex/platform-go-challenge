package utils

import (
	"context"
	"log"
	"time"
)

// Retry function with exponential backoff
func RetryWithExponentialBackoff(ctx context.Context, operation func() error) error {
	const maxAttempts = 3
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		backoff := time.Duration(attempt) * time.Second

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}
	log.Println("Reached all the retry attemps ")
	return err
}
