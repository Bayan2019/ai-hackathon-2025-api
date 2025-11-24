-- +goose Up
CREATE TABLE codes(
    email TEXT NOT NULL REFERENCES users(email) ON DELETE CASCADE,
    -- role roles NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    code TEXT NOT NULL DEFAULT '111111',
    -- objective objectives NOT NULL,
    confirmed TEXT NOT NULL DEFAULT 'FALSE',
    PRIMARY KEY(email, created_at)
);

-- +goose Down
DROP TABLE codes;