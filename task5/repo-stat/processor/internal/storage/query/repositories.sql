-- name: GetRepository :one
SELECT owner, repo, full_name, description, stars, forks, created_at, status, error, updated_at
FROM repositories
WHERE owner = $1 AND repo = $2;

-- name: UpsertRepositoryPending :exec
INSERT INTO repositories (owner, repo, status, error, updated_at)
VALUES ($1, $2, 'pending', '', now())
ON CONFLICT (owner, repo) DO UPDATE
SET status = 'pending',
    error = '',
    updated_at = now();

-- name: UpsertRepositoryInfo :exec
INSERT INTO repositories (owner, repo, full_name, description, stars, forks, created_at, status, error, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'ready', '', now())
ON CONFLICT (owner, repo) DO UPDATE
SET full_name = EXCLUDED.full_name,
    description = EXCLUDED.description,
    stars = EXCLUDED.stars,
    forks = EXCLUDED.forks,
    created_at = EXCLUDED.created_at,
    status = 'ready',
    error = '',
    updated_at = now();

-- name: UpsertRepositoryError :exec
INSERT INTO repositories (owner, repo, status, error, updated_at)
VALUES ($1, $2, 'error', $3, now())
ON CONFLICT (owner, repo) DO UPDATE
SET status = 'error',
    error = EXCLUDED.error,
    updated_at = now();
