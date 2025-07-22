-- name: GetUserTokens :many
SELECT * FROM refresh_tokens;

-- name: InsertNewUserToken :one
INSERT INTO refresh_tokens (
  token,
  user_id,
  expires_at
) VALUES (
    $1, $2, $3
) RETURNING token, expires_at, created_at;

