package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"log"
	"os"
)

func NewPgxPool(ctx context.Context, connStr string, logger tracelog.Logger, logLvl tracelog.LogLevel) (*pgxpool.Pool, error) {

	conf, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, err
	}

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: logLvl,
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		err = fmt.Errorf("Error has occured while creating pool ::: %v", err.Error())
		return nil, err
	}
	return pool, nil
}

func LogLevelFromConfig() (tracelog.LogLevel, error) {

	if stri := os.Getenv("PGX_LOG_LEVEL"); stri != "" {
		logLevel, err := tracelog.LogLevelFromString(stri)
		if err != nil {
			return tracelog.LogLevelDebug, fmt.Errorf("erro while setting loglevel ::: %v", err)

		}
		return logLevel, nil
	}

	return tracelog.LogLevelDebug, nil
}

type PGXStdLogger struct{}

func (p *PGXStdLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	args := make([]any, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

func PgErrors(err error) error {
	var pgEr *pgconn.PgError
	if !errors.As(err, &pgEr) {
		return err
	}
	return fmt.Errorf(

		`Code ::: %v /n,
	Detail ::: %v /n,
	Hint ::: %v /n,
	Position ::: %v /n,
	InternalPosition ::: %v /n,
	InternalQuery ::: %v /n,
	Where ::: %v /n,
	SchemaName ::: %v /n,
	TableName ::: %v /n,
	ColumnName ::: %v /n,
	DataTypeName ::: %v /n,
	ConstraintName ::: %v /n,
	File ::: %v /n,
	Routine ::: %v /n,`,
		pgEr.Code,
		pgEr.Detail,
		pgEr.Hint,
		pgEr.Position,
		pgEr.InternalPosition,
		pgEr.InternalQuery,
		pgEr.Where,
		pgEr.SchemaName,
		pgEr.TableName,
		pgEr.ColumnName,
		pgEr.DataTypeName,
		pgEr.ConstraintName,
		pgEr.File,
		pgEr.Routine,
	)
}
