// Package httputil provides utils for http
package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Get is http.Get, but with context
func Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Post is http.Post, but with context
func Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}

// DoFor200 does a http.Client.Do(), but with a 200 check and error checks
func DoFor200(client *http.Client, request *http.Request) ([]byte, error) {
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // failing is ok

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = nil
		}
		return nil, fmt.Errorf("returns a non-200 status code: %v (%v) with body: %v",
			resp.StatusCode, resp.Status, string(body))
	}
	return io.ReadAll(resp.Body)
}
