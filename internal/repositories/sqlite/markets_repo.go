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

// Insert implements repositories.MarketRepo.
func (m marketRepoImp) Insert(_ context.Context, market pkg.Market) (int, error) {
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

func NewMarketRepo(connString string) (repositories.MarketRepo, error) {

	db, err := getConn(connString)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open db")
	}

	_, err = db.Exec(`
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