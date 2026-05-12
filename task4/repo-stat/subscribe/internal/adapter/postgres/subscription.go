package postgres

import (
	"context"
	"errors"
	"fmt"

	"repo-stat/subscribe/internal/domain"
	"repo-stat/subscribe/internal/storage/db"
	"repo-stat/subscribe/internal/usecase"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type SubscriptionRepository struct {
	queries *db.Queries
}

func NewSubscriptionRepository(queries *db.Queries) *SubscriptionRepository {
	return &SubscriptionRepository{
		queries: queries,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, owner string, repo string) (domain.Subscription, error) {
	subscription, err := r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Subscription{}, usecase.ErrSubscriptionExists
		}
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}

	return domain.Subscription{
		Owner: subscription.Owner,
		Repo:  subscription.Repo,
	}, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, owner string, repo string) error {
	rowsAffected, err := r.queries.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if rowsAffected == 0 {
		return usecase.ErrSubscriptionNotFound
	}

	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]domain.Subscription, error) {
	items, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	subscriptions := make([]domain.Subscription, 0, len(items))
	for _, item := range items {
		subscriptions = append(subscriptions, domain.Subscription{
			Owner: item.Owner,
			Repo:  item.Repo,
		})
	}

	return subscriptions, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
