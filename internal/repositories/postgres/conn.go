package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
)

func ConnectPostgres(ctx context.Context, config configs.SectionPostgres) (*pgxpool.Pool, error) {

	cfg, _ := pgxpool.ParseConfig("")
	cfg.ConnConfig.Host = config.Host
	cfg.ConnConfig.Port = config.Port
	cfg.ConnConfig.User = config.User
	cfg.ConnConfig.Password = config.Pass
	cfg.ConnConfig.Database = config.Database

	cfg.MaxConnIdleTime = config.ConnMaxIdleTime
	cfg.MaxConnLifetime = config.ConnMaxTime
	cfg.MinConns = config.MinConn
	cfg.MaxConns = config.MaxConn

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't make pool")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "couldn't ping database")
	}
	return pool, nil

}
