package sqlite

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type marketStatsImp struct {
	db *sql.DB
}

func (m marketStatsImp) migrate(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS market_stats (
			market_id BIGINTEGER AUTO INCREMENT,
			buy_price TEXT,
			sell_price TEXT,
			date TIMESTAMP,
			PRIMARY KEY(market_id, date)
		)
		`)
	return err
}

func (m marketStatsImp) Insert(ctx context.Context, stat pkg.MarketStat) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO market_stats(market_id, buy_price, sell_price, date) VALUES (?, ?, ?, ?)
	`, stat.MarketID, stat.BuyPrice, stat.SellPrice, stat.Date)

	if err != nil {
		return errors.Wrap(err, "couldn't insert into markets")
	}
	return nil

}

func NewMarketStatsRepo(ctx context.Context, db *sql.DB) (repositories.MarketStatsRepo, error) {

	ans := marketStatsImp{db: db}
	if err := ans.migrate(ctx); err != nil {
		return nil, errors.Wrap(err, "couldn't do the migration")
	}
	return ans, nil

}
