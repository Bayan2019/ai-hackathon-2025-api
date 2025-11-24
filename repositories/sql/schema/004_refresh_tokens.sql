-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email TEXT NOT NULL REFERENCES users(email) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    revoked_at TEXT
);

-- +goose Down
DROP TABLE refresh_tokens;