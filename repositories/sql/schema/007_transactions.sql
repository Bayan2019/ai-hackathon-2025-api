-- +goose Up
CREATE TABLE transactions(
    cst_dim_id INTEGER,
    transdatetime TEXT NOT NULL DEFAULT '',
    transdate TEXT NOT NULL DEFAULT '',
    amount INTEGER NOT NULL DEFAULT 0,
    direction TEXT NOT NULL DEFAULT '',
    target NTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE transactions;