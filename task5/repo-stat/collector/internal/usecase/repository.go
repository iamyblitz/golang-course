package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"repo-stat/collector/internal/domain"
)

var ErrRepositoryNotFound = errors.New("repository not found")

type Repository struct {
	client        *http.Client
	subscriptions SubscriptionProvider
}

type SubscriptionProvider interface {
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

func NewRepository(subscriptions SubscriptionProvider) *Repository {
	return &Repository{
		client: &http.Client{
			Timeout: 12 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
		subscriptions: subscriptions,
	}
}

type githubRepositoryResponse struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int64  `json:"stargazers_count"`
	ForksCount      int64  `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
}

func (u *Repository) GetInfo(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return domain.RepositoryInfo{}, fmt.Errorf("build github request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "repo-stat-collector")

	resp, err := u.do(req)
	if err != nil {
		return domain.RepositoryInfo{}, fmt.Errorf("request github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.RepositoryInfo{}, ErrRepositoryNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.RepositoryInfo{}, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var githubRepo githubRepositoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&githubRepo); err != nil {
		return domain.RepositoryInfo{}, fmt.Errorf("decode github response: %w", err)
	}

	fullName := githubRepo.FullName
	if fullName == "" {
		fullName = owner + "/" + repo
	}

	return domain.RepositoryInfo{
		FullName:    fullName,
		Description: githubRepo.Description,
		Stars:       githubRepo.StargazersCount,
		Forks:       githubRepo.ForksCount,
		CreatedAt:   githubRepo.CreatedAt,
	}, nil
}

func (u *Repository) GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error) {
	subscriptions, err := u.subscriptions.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	repositories := make([]domain.RepositoryInfo, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		info, err := u.GetInfo(ctx, subscription.Owner, subscription.Repo)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, info)
	}

	return repositories, nil
}

func (u *Repository) do(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt < 2; attempt++ {
		resp, err := u.client.Do(req.Clone(req.Context()))
		if err == nil {
			return resp, nil
		}

		lastErr = err
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(200 * time.Millisecond):
		}
	}

	return nil, lastErr
}
