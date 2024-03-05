package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/url"
)

func NewConnection(ctx context.Context, cfg *Config, logger *slog.Logger) (*pgxpool.Pool, error) {

	// build connection string from config
	dsn := url.URL{
		Scheme: "postgres",
		Host:   cfg.Host,
		User:   url.UserPassword(cfg.User, cfg.Password),
		Path:   cfg.Dbname,
	}
	q := dsn.Query()
	q.Add("sslmode", cfg.SSLMode)
	dsn.RawQuery = q.Encode()

	logger.Info("connecting to psql", "host", cfg.Host, "database", cfg.Dbname)

	// create pgx config
	pgxConfig, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		logger.Error("error while parsing psql config", err)
		return nil, err
	}

	// register custom data type handlers
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		RegisterTypeUUID(conn.TypeMap())
		return nil
	}

	// connect to pg
	dbpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		logger.Error("error while connecting to psql", err)
		return nil, err
	}

	// ensure connection is established
	if err = dbpool.Ping(ctx); err != nil {
		logger.Error("error while pinging psql", err)
		return nil, err
	}

	logger.Info("successfully connected to psql")
	return dbpool, nil
}
