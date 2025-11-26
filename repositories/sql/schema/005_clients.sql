-- +goose Up
CREATE TABLE clients(
    cst_dim_id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    gender TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE clients;