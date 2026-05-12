package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type RepositoryProvider interface {
	GetRepositoryInfo(ctx context.Context, url string) (domain.RepositoryInfo, error)
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

func (u *Repository) GetInfo(ctx context.Context, url string) (domain.RepositoryInfo, error) {
	return u.provider.GetRepositoryInfo(ctx, url)
}

func (u *Repository) GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error) {
	return u.provider.GetSubscriptionsInfo(ctx)
}
