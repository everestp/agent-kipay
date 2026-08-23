package database

import (
	"context"

	"github.com/everest/bheri/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() (*pgxpool.Pool, error) {
	cfg := config.LoadDatabaseConfig()

	poolConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		return nil, err
	}

	// Avoid pgx prepared-statement caching issues.
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
