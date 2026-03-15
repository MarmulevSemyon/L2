package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	client *http.Client
}

type Response struct {
	FinalURL    string
	StatusCode  int
	ContentType string
	Body        []byte
}

func New(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (f *Fetcher) Fetch(rawURL string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return Response{}, fmt.Errorf("create request: %w", err)
	}

	// req.Header.Set("User-Agent", "mini-wget/1.0")

	resp, err := f.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("read response body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")

	return Response{
		FinalURL:    resp.Request.URL.String(),
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Body:        body,
	}, nil
}
