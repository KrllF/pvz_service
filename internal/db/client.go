package db

import (
	"context"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewDB конструктор Database
func NewDB(ctx context.Context, conf config.Config) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, conf.ConnectConfig.DSN)
	if err != nil {
		return nil, err
	}

	return NewDatabase(pool), nil
}
