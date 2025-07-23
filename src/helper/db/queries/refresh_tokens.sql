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

-- name: GetToken :one
SELECT token, user_id FROM refresh_tokens
WHERE token = $1
AND expires_at > now()
LIMIT 1;

-- name: DeleteToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteExpiresToken :exec
DELETE FROM refresh_tokens 
WHERE expires_at <= now();

