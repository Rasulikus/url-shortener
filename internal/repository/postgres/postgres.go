package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasulikus/url-shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	poolcfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	poolcfg.MinConns = int32(cfg.MinConns)
	poolcfg.MaxConns = int32(cfg.MaxConns)
	poolcfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolcfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	poolcfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
	defer pingCancel()

	if err = pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
