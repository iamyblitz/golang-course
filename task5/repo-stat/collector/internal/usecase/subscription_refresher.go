package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type TaskPublisher interface {
	PublishCollectTask(ctx context.Context, owner string, repo string) error
}

type SubscriptionRefresher struct {
	log           *slog.Logger
	subscriptions SubscriptionProvider
	publisher     TaskPublisher
	interval      time.Duration
}

func NewSubscriptionRefresher(
	log *slog.Logger,
	subscriptions SubscriptionProvider,
	publisher TaskPublisher,
	interval time.Duration,
) *SubscriptionRefresher {
	return &SubscriptionRefresher{
		log:           log,
		subscriptions: subscriptions,
		publisher:     publisher,
		interval:      interval,
	}
}

func (r *SubscriptionRefresher) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := r.refresh(ctx); err != nil {
				r.log.Error("subscriptions refresh failed", "error", err)
			}
		}
	}
}

func (r *SubscriptionRefresher) refresh(ctx context.Context) error {
	subscriptions, err := r.subscriptions.ListSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("list subscriptions: %w", err)
	}

	for _, subscription := range subscriptions {
		if err := r.publisher.PublishCollectTask(ctx, subscription.Owner, subscription.Repo); err != nil {
			return fmt.Errorf("publish subscription refresh task: %w", err)
		}
		r.log.Debug("published subscription refresh task", "owner", subscription.Owner, "repo", subscription.Repo)
	}

	return nil
}
