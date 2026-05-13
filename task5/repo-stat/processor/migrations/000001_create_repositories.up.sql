CREATE TABLE repositories (
    owner TEXT NOT NULL,
    repo TEXT NOT NULL,
    full_name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    stars BIGINT NOT NULL DEFAULT 0,
    forks BIGINT NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    error TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (owner, repo)
);
