package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// AssetsFetcher is a worker pool to download assets from the assets server.
type AssetsFetcher struct {
	cli *http.Client
	srv *url.URL
}

// FetchResult is the asset server GET response for a key
type FetchResult struct {
	Key  string
	Body []byte
	Err  error
}

// newFetchers starts n workers to fetch assets.
// It returns two channels (in,out). It accepts keys from in and returns responses to out.
// It runs until in is closed. When finished, it closes out.
func (af *AssetsFetcher) newFetcher(ctx context.Context, n int) (chan<- string, <-chan *FetchResult) {
	in := make(chan string)
	out := make(chan *FetchResult, n*2)
	// out must be buffered so that fetchers do not have to wait the consumer
	// Because the consumer is very fast, n*2 should do the job

	go func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				for key := range in {
					out <- af.fetch(ctx, key)
				}
			}()
		}
		wg.Wait()
		close(out)
	}(ctx)

	return in, out
}

// fetch is a simple GET client. Calls the asset server for key and returns the response.
// Expects only 200 OK responses with size than 1MB.
func (af *AssetsFetcher) fetch(ctx context.Context, key string) *FetchResult {
	uri := af.srv.JoinPath(key)
	req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		return &FetchResult{Key: key, Err: fmt.Errorf("fetcher: failed to create request: %w", err)}
	}

	token := ctx.Value(keyToken).(string)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := af.cli.Do(req)
	if err != nil {
		return &FetchResult{Key: key, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &FetchResult{Key: key, Err: fmt.Errorf("fetcher: got status code %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(&io.LimitedReader{R: resp.Body, N: 1048576})
	if err != nil {
		return &FetchResult{Key: key, Err: fmt.Errorf("fetcher: read failed: %w", err)}
	}
	return &FetchResult{Key: key, Body: body}
}
