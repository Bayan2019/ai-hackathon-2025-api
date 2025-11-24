-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens(token, email, expires_at)
VALUES (?, ?, ?);
--

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.email = refresh_tokens.email
WHERE refresh_tokens.token = ?
    AND revoked_at IS NULL
    AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;
--

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = CURRENT_TIMESTAMP, revoked_at = CURRENT_TIMESTAMP
WHERE token = ? AND revoked_at IS NULL;
--

-- name: GetRefreshTokenOfUser :one
SELECT token FROM refresh_tokens
WHERE email = ?
    AND revoked_at IS NULL
    AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;
--