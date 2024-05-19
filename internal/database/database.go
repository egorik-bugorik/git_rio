package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func NewPgxPool(ctx context.Context, connStr string, logger tracelog.Logger, logLvl tracelog.LogLevel) (*pgxpool.Pool, error) {

	return nil, nil
}
