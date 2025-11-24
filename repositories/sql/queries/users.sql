-- name: CreateUser :exec
INSERT INTO users(email, password_hash, first_name, last_name)
VALUES (?, ?, ?, ?);
--

-- name: GetUsers :many
SELECT * FROM users;
--

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;
--