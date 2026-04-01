package db

import (
	"context"
	"fmt"
	"time"

	"restaurant-management-system/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is a global PostgreSQL client connection pool.
var Client *pgxpool.Pool

// InitPostgres initializes the global PostgreSQL client using pgx.
func InitPostgres() error {
	dbCfg := config.RuntimeConfig.Database

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",

		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to create postgres pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	Client = pool
	return nil
}

// ClosePostgres closes the global PostgreSQL client.
func ClosePostgres() {
	if Client != nil {
		Client.Close()
	}
}
