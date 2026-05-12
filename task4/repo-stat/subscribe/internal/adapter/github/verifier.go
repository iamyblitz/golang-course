package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"repo-stat/subscribe/internal/usecase"
)

type Verifier struct {
	client *http.Client
}

func NewVerifier() *Verifier {
	return &Verifier{
		client: &http.Client{
			Timeout: 12 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
	}
}

func (v *Verifier) Exists(ctx context.Context, owner string, repo string) error {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("build github request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "repo-stat-subscribe")

	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: request github api: %v", usecase.ErrRepositoryUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return usecase.ErrRepositoryNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: github api returned status %d", usecase.ErrRepositoryUnavailable, resp.StatusCode)
	}

	return nil
}
