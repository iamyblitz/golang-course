package postgres

import (
	"context"
	"errors"
	"fmt"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/storage/db"
	"repo-stat/processor/internal/usecase"

	"github.com/jackc/pgx/v5"
)

type RepositoryStore struct {
	queries *db.Queries
}

func NewRepositoryStore(queries *db.Queries) *RepositoryStore {
	return &RepositoryStore{queries: queries}
}

func (s *RepositoryStore) Get(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, string, string, error) {
	item, err := s.queries.GetRepository(ctx, db.GetRepositoryParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RepositoryInfo{}, "", "", usecase.ErrRepositoryNotCached
		}
		return domain.RepositoryInfo{}, "", "", fmt.Errorf("get repository: %w", err)
	}

	return domain.RepositoryInfo{
		FullName:    item.FullName,
		Description: item.Description,
		Stars:       item.Stars,
		Forks:       item.Forks,
		CreatedAt:   item.CreatedAt,
	}, item.Status, item.Error, nil
}

func (s *RepositoryStore) MarkPending(ctx context.Context, owner string, repo string) error {
	return s.queries.UpsertRepositoryPending(ctx, db.UpsertRepositoryPendingParams{
		Owner: owner,
		Repo:  repo,
	})
}

func (s *RepositoryStore) SaveInfo(ctx context.Context, owner string, repo string, info domain.RepositoryInfo) error {
	return s.queries.UpsertRepositoryInfo(ctx, db.UpsertRepositoryInfoParams{
		Owner:       owner,
		Repo:        repo,
		FullName:    info.FullName,
		Description: info.Description,
		Stars:       info.Stars,
		Forks:       info.Forks,
		CreatedAt:   info.CreatedAt,
	})
}

func (s *RepositoryStore) SaveError(ctx context.Context, owner string, repo string, message string) error {
	return s.queries.UpsertRepositoryError(ctx, db.UpsertRepositoryErrorParams{
		Owner: owner,
		Repo:  repo,
		Error: message,
	})
}
