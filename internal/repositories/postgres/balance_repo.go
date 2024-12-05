package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/repositories"
)

type BalanceRepoImp struct {
	pool *pgxpool.Pool
}

func (m BalanceRepoImp) migrate(ctx context.Context) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS balances(
			exchange_name varchar(50),
			date TIMESTAMP,
			balance TEXT,
			PRIMARY KEY(exchange_name, date)
		)
		`
	_, err := m.pool.Exec(ctx, stmt)
	return err

}

func (m BalanceRepoImp) Insert(ctx context.Context, exchangeName string, date time.Time, balance decimal.Decimal) error {
	stmt := `
		INSERT INTO balances(
			exchange_name,
			date,
			balance
		) VALUES (
			$1, $2, $3
		)
		`
	_, err := m.pool.Exec(ctx, stmt,
		exchangeName,
		date,
		balance,
	)

	return err

}

func (m BalanceRepoImp) Get(ctx context.Context, exchangeName string, start, end time.Time) ([]time.Time, []decimal.Decimal, error) {
	stmt := `SELECT date, balance FROM balances WHERE exchange_name = $1`

	rows, err := m.pool.Query(ctx, stmt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get rows")
	}

	var dates []time.Time
	var balances []decimal.Decimal

	for rows.Next() {
		var date time.Time
		var balance decimal.Decimal

		if err = rows.Scan(&date, &balance); err != nil {
			return nil, nil, errors.Wrap(err, "couldn't scan to market")
		}
		dates = append(dates, date)
		balances = append(balances, balance)
	}
	return dates, balances, nil
}

func NewBalanceRepo(pool *pgxpool.Pool, ctx context.Context) (repositories.BalanceRepo, error) {
	ans := BalanceRepoImp{
		pool: pool,
	}

	return ans, ans.migrate(ctx)
}
