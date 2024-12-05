package sqlite

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type marketRepoImp struct {
	db *sql.DB
}

func (m marketRepoImp) migrate(ctx context.Context) error {

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS markets (
			id INTEGER PRIMARY KEY,
			exchange_name TEXT,
			base_asset TEXT,
			quote_asset TEXT,
			UNIQUE(exchange_name, base_asset, quote_asset)
		)`,
		`ALTER TABLE markets ADD fa_name TEXT`,
		`ALTER TABLE markets ADD en_name TEXT`,
		`ALTER TABLE markets ADD  is_active BOOL DEFAULT FALSE`,
	}
	for _, stmt := range stmts {
		if _, err := m.db.ExecContext(ctx, stmt); err != nil {
			return errors.Wrap(errors.WithMessage(err, "statements is: "+stmt), "couldn't do the migration")
		}
	}
	return nil
}
func (m marketRepoImp) GetOrCreate(ctx context.Context, market pkg.Market) (int, bool, error) {
	_, err := m.db.ExecContext(ctx, `
INSERT INTO markets(exchange_name, base_asset, quote_asset, en_name, fa_name) VALUES (?, ?, ?, ?, ?) ON CONFLICT DO NOTHING
	`, market.ExchangeName, market.BaseAsset, market.QuoteAsset, market.EnName, market.FaName)
	if err != nil {
		return -1, false, errors.Wrap(err, "couldn't insert into markets")
	}
	var marketID int
	var isActive bool
	err = m.db.QueryRow("SELECT id, is_active FROM markets WHERE exchange_name = ? AND base_asset = ? AND quote_asset = ?", market.ExchangeName, market.BaseAsset, market.QuoteAsset).Scan(&marketID, &isActive)
	if err != nil {
		return -1, false, errors.Wrap(err, "couldn't get id ")
	}
	return marketID, isActive, nil

}

func (m marketRepoImp) GetAllExchangeMarkets(ctx context.Context, exchangeName string) ([]pkg.Market, error) {
	stmt := `SELECT id, base_asset, quote_asset FROM markets WHERE exchange_name = ?`
	rows, err := m.db.QueryContext(ctx, stmt, exchangeName)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get rows")
	}

	var markets []pkg.Market
	for rows.Next() {
		var market pkg.Market
		if err = rows.Scan(&market.ID, &market.BaseAsset, &market.QuoteAsset); err != nil {
			return nil, errors.Wrap(err, "couldn't scan to market")
		}
		market.ExchangeName = exchangeName
		markets = append(markets, market)
	}
	return markets, nil
}

func (m marketRepoImp) GetByID(ctx context.Context, marketID int) (pkg.Market, error) {
	stmt := `SELECT exchange_name, base_asset, quote_asset FROM markets WHERE id = ?`
	var ans pkg.Market
	ans.ID = marketID

	err := m.db.QueryRowContext(ctx, stmt, marketID).Scan(&ans.ExchangeName, &ans.BaseAsset, &ans.QuoteAsset)
	if err != nil {
		return ans, errors.Wrap(err, "couldn't get data from db")
	}

	return ans, nil
}

func NewMarketRepo(ctx context.Context, db *sql.DB) (repositories.MarketRepo, error) {
	ans := marketRepoImp{db: db}
	return ans, ans.migrate(ctx)
}
