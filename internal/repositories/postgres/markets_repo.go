package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type MarketsRepoImp struct {
	pool *pgxpool.Pool
}

func (m MarketsRepoImp) migrate(ctx context.Context) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS markets(
			id SERIAL PRIMARY KEY,
			exchange_name varchar(50),
			base_asset varchar(50),
			quote_asset varchar(50),
			UNIQUE(exchange_name, base_asset, quote_asset)
		)
		`
	_, err := m.pool.Exec(ctx, stmt)
	return err

}

func (m MarketsRepoImp) GetOrCreate(ctx context.Context, market pkg.Market) (int, error) {
	stmt := `
		INSERT INTO markets(
			exchange_name,
			base_asset,
			quote_asset
		) VALUES (
			$1, $2, $3
		) ON CONFLICT(exchange_name, base_asset, quote_asset) DO UPDATE SET id = markets.id
			RETURNING id
		`
	var marketID int
	err := m.pool.QueryRow(ctx, stmt,
		market.ExchangeName,
		market.BaseAsset,
		market.QuoteAsset,
	).Scan(&marketID)

	return marketID, err

}

func (m MarketsRepoImp) GetAllExchangeMarkets(ctx context.Context, exchangeName string) ([]pkg.Market, error) {
	stmt := `SELECT id, base_asset, quote_asset FROM markets WHERE exchange_name = $1`
	rows, err := m.pool.Query(ctx, stmt, exchangeName)
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

func (m MarketsRepoImp) GetByID(ctx context.Context, marketID int) (pkg.Market, error) {
	stmt := `SELECT exchange_name, base_asset, quote_asset FROM markets WHERE id = $1`
	var ans pkg.Market
	ans.ID = marketID

	err := m.pool.QueryRow(ctx, stmt).Scan(&ans.ExchangeName, &ans.BaseAsset, &ans.QuoteAsset)
	if err != nil {
		return ans, errors.Wrap(err, "couldn't get data from db")
	}

	return ans, nil
}

func NewMarketRepo(pool *pgxpool.Pool, ctx context.Context) (repositories.MarketRepo, error) {
	ans := MarketsRepoImp{
		pool: pool,
	}
	return ans, ans.migrate(ctx)
}
