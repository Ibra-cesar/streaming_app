package helper

import (
	"context"
	"fmt"

	"github.com/Ibra-cesar/video-streaming/src/internal/handlers"
	"github.com/Ibra-cesar/video-streaming/src/internal/query_repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct{
	Pool *pgxpool.Pool
	AuthHandler *handlers.AuthConnServices
	RefreshHandler *handlers.RefreshHandlers
	Queries *query_repo.Queries
}

func App(ctx context.Context) (*Application, error){
	pool, err := DbConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to Connect to DataBase %w", err)
	}

	queries := query_repo.New(pool)


	jwtSecrets := []byte(Env("JWT_SECRET_KEY"))
	if len(jwtSecrets) == 0 {
		jwtSecrets = []byte("JWT_FAKE_KEY")
	}
	jwtRefreshTokenSecret := []byte(Env("JWT_REFRESH_TOKEN_KEY"))

	authHandler := handlers.AuthServices(queries, jwtSecrets, jwtRefreshTokenSecret)
	refreshHandler := handlers.RefreshService(*queries, jwtSecrets, jwtRefreshTokenSecret)

	return &Application{
		Pool: pool,
		AuthHandler: authHandler,
		RefreshHandler: refreshHandler,
		Queries: queries,
	}, nil
}
