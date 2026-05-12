CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    owner TEXT NOT NULL,
    repo TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_owner_repo_unique UNIQUE (owner, repo)
);
