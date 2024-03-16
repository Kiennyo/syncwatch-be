package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kiennyo/syncwatch-be/internal/config"
)

func New(ctx context.Context, cfg config.DB) (*pgxpool.Pool, error) {
	parseConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	parseConfig.MaxConns = int32(cfg.MaxOpenConn)
	parseConfig.MinConns = int32(cfg.MaxIdleConn)
	parseConfig.MaxConnIdleTime = duration

	db, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("database connection pool established")

	return db, nil
}
