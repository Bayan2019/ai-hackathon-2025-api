-- +goose Up
CREATE TABLE transactions(
    cst_dim_id INTEGER NOT NULL REFERENCES clients(cst_dim_id) ON DELETE CASCADE,
    transdatetime TEXT NOT NULL DEFAULT '',
    transdate TEXT NOT NULL DEFAULT '',
    amount INTEGER NOT NULL DEFAULT 0,
    direction TEXT NOT NULL DEFAULT '',
    target INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE transactions;