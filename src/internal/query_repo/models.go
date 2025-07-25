// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package query_repo

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt pgtype.Timestamptz
	CreatedAt pgtype.Timestamptz
}

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	IsAdmin      pgtype.Bool
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}
