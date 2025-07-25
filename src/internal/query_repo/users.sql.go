// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package query_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
  WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getAllPlayers = `-- name: GetAllPlayers :many
SELECT id, name, email, password_hash, is_admin, created_at, updated_at FROM users
`

func (q *Queries) GetAllPlayers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllPlayers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.PasswordHash,
			&i.IsAdmin,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, email, name FROM users where id = $1
`

type GetUserRow struct {
	ID    uuid.UUID
	Email string
	Name  string
}

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (GetUserRow, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(&i.ID, &i.Email, &i.Name)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT password_hash, id, email, name FROM users WHERE email = $1
`

type GetUserByEmailRow struct {
	PasswordHash string
	ID           uuid.UUID
	Email        string
	Name         string
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.PasswordHash,
		&i.ID,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one

INSERT INTO users (
  id,
  name,
  email,
  password_hash
) VALUES ( 
    $1, $2, $3, $4
) RETURNING id, name, email, is_admin, created_at, updated_at
`

type InsertUserParams struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
}

type InsertUserRow struct {
	ID        uuid.UUID
	Name      string
	Email     string
	IsAdmin   pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (InsertUserRow, error) {
	row := q.db.QueryRow(ctx, insertUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
	)
	var i InsertUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
