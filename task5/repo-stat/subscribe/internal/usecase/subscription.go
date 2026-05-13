package usecase

import (
	"context"
	"errors"

	"repo-stat/subscribe/internal/domain"
)

var (
	ErrInvalidSubscription   = errors.New("owner and repo are required")
	ErrSubscriptionExists    = errors.New("subscription already exists")
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrRepositoryNotFound    = errors.New("repository not found")
	ErrRepositoryUnavailable = errors.New("repository verifier unavailable")
)

type SubscriptionRepository interface {
	Create(ctx context.Context, owner string, repo string) (domain.Subscription, error)
	Delete(ctx context.Context, owner string, repo string) error
	List(ctx context.Context) ([]domain.Subscription, error)
}

type RepositoryVerifier interface {
	Exists(ctx context.Context, owner string, repo string) error
}

type Subscription struct {
	repository SubscriptionRepository
	verifier   RepositoryVerifier
}

func NewSubscription(repository SubscriptionRepository, verifier RepositoryVerifier) *Subscription {
	return &Subscription{
		repository: repository,
		verifier:   verifier,
	}
}

func (u *Subscription) Create(ctx context.Context, owner string, repo string) (domain.Subscription, error) {
	if owner == "" || repo == "" {
		return domain.Subscription{}, ErrInvalidSubscription
	}

	if err := u.verifier.Exists(ctx, owner, repo); err != nil {
		return domain.Subscription{}, err
	}

	return u.repository.Create(ctx, owner, repo)
}

func (u *Subscription) Delete(ctx context.Context, owner string, repo string) error {
	if owner == "" || repo == "" {
		return ErrInvalidSubscription
	}

	return u.repository.Delete(ctx, owner, repo)
}

func (u *Subscription) List(ctx context.Context) ([]domain.Subscription, error) {
	return u.repository.List(ctx)
}
