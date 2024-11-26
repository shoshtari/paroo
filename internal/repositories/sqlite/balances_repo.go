package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/repositories"
)

type BalanceRepoImp struct {
	db *sql.DB
}

func (m BalanceRepoImp) migrate(ctx context.Context) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS balances(
			exchange_name TEXT,
			date TIMESTAMP,
			balance TEXT,
			PRIMARY KEY(exchange_name, date)
		)
		`
	_, err := m.db.ExecContext(ctx, stmt)
	return err

}

func (m BalanceRepoImp) Insert(ctx context.Context, exchangeName string, date time.Time, balance decimal.Decimal) error {
	stmt := `
		INSERT INTO balances(
			exchange_name,
			date,
			balance
		) VALUES (
			?, ?, ?
		)
		`
	_, err := m.db.ExecContext(ctx, stmt,
		exchangeName,
		date,
		balance.String(),
	)

	if err != nil {
		return errors.WithStack(err)
	}
	return nil

}

func (m BalanceRepoImp) Get(ctx context.Context, exchangeName string, start, end time.Time) ([]time.Time, []decimal.Decimal, error) {
	stmt := `SELECT date, balance FROM balances WHERE exchange_name = ?`

	rows, err := m.db.QueryContext(ctx, stmt)
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

func NewBalanceRepo(ctx context.Context, db *sql.DB) (repositories.BalanceRepo, error) {
	ans := BalanceRepoImp{
		db: db,
	}

	return ans, ans.migrate(ctx)
}
