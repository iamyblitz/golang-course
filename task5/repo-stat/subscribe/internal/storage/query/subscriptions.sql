-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING id, owner, repo, created_at;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;

-- name: ListSubscriptions :many
SELECT id, owner, repo, created_at
FROM subscriptions
ORDER BY owner, repo;
