package usecase

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"repo-stat/processor/internal/domain"
)

var ErrInvalidRepositoryURL = errors.New("invalid repository url")

type RepositoryProvider interface {
	GetRepositoryInfo(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error)
	GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error)
}

type Repository struct {
	provider RepositoryProvider
}

func NewRepository(provider RepositoryProvider) *Repository {
	return &Repository{
		provider: provider,
	}
}

func (u *Repository) GetInfo(ctx context.Context, repoURL string) (domain.RepositoryInfo, error) {
	owner, repo, err := parseGitHubRepositoryURL(repoURL)
	if err != nil {
		return domain.RepositoryInfo{}, err
	}

	return u.provider.GetRepositoryInfo(ctx, owner, repo)
}

func (u *Repository) GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error) {
	return u.provider.GetSubscriptionsInfo(ctx)
}

func parseGitHubRepositoryURL(repoURL string) (string, string, error) {
	parsedURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return "", "", ErrInvalidRepositoryURL
	}

	if parsedURL.Scheme != "https" || parsedURL.Host != "github.com" {
		return "", "", ErrInvalidRepositoryURL
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", ErrInvalidRepositoryURL
	}

	return parts[0], parts[1], nil
}
