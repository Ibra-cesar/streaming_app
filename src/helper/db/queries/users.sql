-- name: GetAllPlayers :many
SELECT * FROM users;

-- name: InsertUser :one

INSERT INTO users (
  id,
  name,
  email,
  password_hash
) VALUES ( 
    $1, $2, $3, $4
) RETURNING id, name, email, is_admin, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
  WHERE id = $1;

-- name: GetUser :one
SELECT id, email FROM users where id = $1;
