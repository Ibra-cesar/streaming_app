package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolConfig struct{
	URI string
	MaxConn int32
	MinConn int32
	MaxLifeTime time.Duration
	MaxIdleTime time.Duration
}

func DefaultConfig() *PoolConfig{
	env := Env("DB_URI")
	return &PoolConfig{
		URI: env,
		MaxConn: 15,
		MinConn: 5,
		MaxLifeTime: 15 * time.Minute,
		MaxIdleTime: 5 * time.Minute,
	} 
}

func ConnPool(ctx context.Context, conf *PoolConfig) (*pgxpool.Pool, error){
	config, err := pgxpool.ParseConfig(conf.URI)
	if err != nil {
		return nil, fmt.Errorf("Error Parsing Config, %w", err)
	}

	config.MaxConns = conf.MaxConn
	config.MinConns = conf.MinConn
	config.MaxConnLifetime = conf.MaxLifeTime
	config.MaxConnIdleTime = conf.MaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return  nil, fmt.Errorf("Failed to establish connection to database, %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return  nil, fmt.Errorf("Failed to Ping Database, %w", err)
	}

	fmt.Println("Connection established")
	return  pool, nil
}
