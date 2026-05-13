package usecase

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"repo-stat/processor/internal/domain"
)

var (
	ErrInvalidRepositoryURL  = errors.New("invalid repository url")
	ErrRepositoryNotCached   = errors.New("repository not cached")
	ErrRepositoryNotFound    = errors.New("repository not found")
	ErrRepositoryUnavailable = errors.New("repository unavailable")
)

type RepositoryStore interface {
	Get(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, string, string, error)
	MarkPending(ctx context.Context, owner string, repo string) error
}

type TaskPublisher interface {
	PublishCollectTask(ctx context.Context, owner string, repo string) error
}

type SubscriptionProvider interface {
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type Repository struct {
	store         RepositoryStore
	publisher     TaskPublisher
	subscriptions SubscriptionProvider
	waitTimeout   time.Duration
}

func NewRepository(store RepositoryStore, publisher TaskPublisher, subscriptions SubscriptionProvider) *Repository {
	return &Repository{
		store:         store,
		publisher:     publisher,
		subscriptions: subscriptions,
		waitTimeout:   25 * time.Second,
	}
}

func (u *Repository) GetInfo(ctx context.Context, repoURL string) (domain.RepositoryInfo, error) {
	owner, repo, err := parseGitHubRepositoryURL(repoURL)
	if err != nil {
		return domain.RepositoryInfo{}, err
	}

	return u.getOrRequest(ctx, owner, repo)
}

func (u *Repository) GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error) {
	subscriptions, err := u.subscriptions.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	repositories := make([]domain.RepositoryInfo, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		info, err := u.getOrRequest(ctx, subscription.Owner, subscription.Repo)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, info)
	}

	return repositories, nil
}

func (u *Repository) getOrRequest(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error) {
	info, statusValue, errorMessage, err := u.store.Get(ctx, owner, repo)
	if err == nil {
		switch statusValue {
		case "ready":
			return info, nil
		case "error":
			if errorMessage == "not_found" {
				return domain.RepositoryInfo{}, ErrRepositoryNotFound
			}
		}
	} else if !errors.Is(err, ErrRepositoryNotCached) {
		return domain.RepositoryInfo{}, err
	}

	if err := u.store.MarkPending(ctx, owner, repo); err != nil {
		return domain.RepositoryInfo{}, err
	}
	if err := u.publisher.PublishCollectTask(ctx, owner, repo); err != nil {
		return domain.RepositoryInfo{}, ErrRepositoryUnavailable
	}

	return u.waitForResult(ctx, owner, repo)
}

func (u *Repository) waitForResult(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error) {
	ctxWait, cancel := context.WithTimeout(ctx, u.waitTimeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctxWait.Done():
			return domain.RepositoryInfo{}, ErrRepositoryUnavailable
		case <-ticker.C:
			info, statusValue, errorMessage, err := u.store.Get(ctxWait, owner, repo)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return domain.RepositoryInfo{}, ErrRepositoryUnavailable
				}
				if errors.Is(err, ErrRepositoryNotCached) {
					continue
				}
				return domain.RepositoryInfo{}, err
			}

			switch statusValue {
			case "ready":
				return info, nil
			case "error":
				return domain.RepositoryInfo{}, repositoryError(errorMessage)
			}
		}
	}
}

func repositoryError(message string) error {
	switch message {
	case "not_found":
		return ErrRepositoryNotFound
	case "github_unavailable":
		return ErrRepositoryUnavailable
	default:
		return ErrRepositoryUnavailable
	}
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
