package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type exchangeRepoImp struct {
	pool *pgxpool.Pool
}

func (e exchangeRepoImp) GetByName(context.Context, string) (int, error) {
	panic("unimplemented")
}

func (e exchangeRepoImp) Insert(ctx context.Context, exchange pkg.Exchange) (int, error) {
	err := e.pool.QueryRow(ctx, `
	INSERT INTO exchanges(name, rial_symbol, tether_symbol) VALUES(
		$1, $2, $3
	), ON CONFLICT(name) DO UPDATE SET name = name
	RETURNING id
	`, exchange.Name, exchange.RialSymbol, exchange.TetherSymbol).Scan(&exchange.ID)

	return exchange.ID, err

}

func (m exchangeRepoImp) migrate(ctx context.Context) error {
	stmts := []string{`
		CREATE TABLE IF NOT EXISTS exchanges(
			id SERIAL PRIMARY KEY,
			name varchar(50) UNIQUE,
			rial_symbol varchar(50),
			tether_symbol varchar(50),
			created_at TIMESTAMP DEFAULT NOW(),
		)
		`,
	}
	for _, stmt := range stmts {
		if _, err := m.pool.Exec(ctx, stmt); err != nil {
			return errors.Wrap(errors.WithMessage(err, "error on stmt: "+stmt), "couldn't do migration")
		}
	}
	return nil

}

func NewExchangeRepo(ctx context.Context, pool *pgxpool.Pool) (repositories.ExchangeRepo, error) {
	ans := exchangeRepoImp{
		pool: pool,
	}
	return ans, ans.migrate(ctx)
}
