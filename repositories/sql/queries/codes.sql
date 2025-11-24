-- name: CreateCode :exec
INSERT INTO codes(email, code)
VALUES (?, ?);
--

-- name: GetCodesOfUser :many
SELECT * FROM codes
WHERE email = ?
ORDER BY created_at DESC;
--

-- name: GetCodeOfUser :one
SELECT * FROM codes
WHERE email = ? AND
    confirmed = 'FALSE' AND
    created_at > DATETIME('now', '-5 minutes')
ORDER BY created_at DESC;
--

-- name: ConfirmCode :exec
UPDATE codes
SET updated_at = CURRENT_TIMESTAMP,
    confirmed = 'TRUE'
WHERE email = ? AND code = ? AND created_at = ?;
--