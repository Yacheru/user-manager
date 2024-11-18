package postgres

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"user-manager/init/logger"
	"user-manager/pkg/constants"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgresConnection(ctx context.Context, dsn string) (*sqlx.DB, error) {
	logger.Debug("new postgresql connection", constants.PostgresCategory)

	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}

	logger.Debug("postgres connected", constants.PostgresCategory)

	m, err := migrate.New("file://./migrations", dsn)
	if err != nil {
		logger.Error(err.Error(), constants.PostgresCategory)
		return nil, err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Debug("migrations already up to date", constants.PostgresCategory)
		} else {
			logger.Error(err.Error(), constants.PostgresCategory)
			return nil, err
		}
	}

	return db, nil
}
