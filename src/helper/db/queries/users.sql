-- name: GetAllPlayers :many
SELECT * FROM users;

-- name: InsertUser :one

INSERT INTO users (
  id,
  name,
  email,
  password_hash,
  salt
) VALUES ( 
    $1, $2, $3, $4, $5
)
RETURNING *;
