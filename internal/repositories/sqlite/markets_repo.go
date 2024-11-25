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

func (m marketRepoImp) GetOrCreate(_ context.Context, market pkg.Market) (int, error) {
	_, err := m.db.Exec(`
INSERT INTO markets(exchange_name, base_asset, quote_asset) VALUES (?, ?, ?) ON CONFLICT DO NOTHING
	`, market.ExchangeName, market.BaseAsset, market.QuoteAsset)
	if err != nil {
		return -1, errors.Wrap(err, "couldn't insert into markets")
	}
	var marketID int
	err = m.db.QueryRow("SELECT id FROM markets WHERE exchange_name = ? AND base_asset = ? AND quote_asset = ?", market.ExchangeName, market.BaseAsset, market.QuoteAsset).Scan(&marketID)
	if err != nil {
		return -1, errors.Wrap(err, "couldn't get id ")
	}
	return marketID, nil

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

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS markets (
			id INTEGER PRIMARY KEY,
			exchange_name TEXT,
			base_asset TEXT,
			quote_asset TEXT,
			UNIQUE(exchange_name, base_asset, quote_asset)
		)
		`)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't do the migration")
	}

	return marketRepoImp{db: db}, nil

}
