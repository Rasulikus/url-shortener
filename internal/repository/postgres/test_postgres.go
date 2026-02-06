package postgres

import (
	"context"
	"time"

	"github.com/Rasulikus/url-shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testDBCfg = &config.DBConfig{
		Host:            "localhost",
		Port:            "54329",
		User:            "test",
		Pass:            "test",
		Name:            "testdb",
		SSLMode:         "disable",
		MinConns:        1,
		MaxConns:        5,
		MaxConnLifetime: 5 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	}
)

func NewTestPool() (*pgxpool.Pool, error) {
	pool, err := NewPool(testDBCfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func TruncateUrls(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `TRUNCATE TABLE urls RESTART IDENTITY;`)
	return err
}
