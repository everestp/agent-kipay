// database/db.go

package database

import (
	"context"
	"fmt"

	"github.com/everest/bheri/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgres() (*pgx.Conn, error) {
	cfg := config.LoadDatabaseConfig()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	return conn, nil
}
