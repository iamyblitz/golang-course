package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type SubscriptionProvider interface {
	CreateSubscription(ctx context.Context, owner string, repo string) (domain.Subscription, error)
	DeleteSubscription(ctx context.Context, owner string, repo string) error
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type Subscription struct {
	provider SubscriptionProvider
}

func NewSubscription(provider SubscriptionProvider) *Subscription {
	return &Subscription{
		provider: provider,
	}
}

func (u *Subscription) Create(ctx context.Context, owner string, repo string) (domain.Subscription, error) {
	return u.provider.CreateSubscription(ctx, owner, repo)
}

func (u *Subscription) Delete(ctx context.Context, owner string, repo string) error {
	return u.provider.DeleteSubscription(ctx, owner, repo)
}

func (u *Subscription) List(ctx context.Context) ([]domain.Subscription, error) {
	return u.provider.ListSubscriptions(ctx)
}
