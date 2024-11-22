package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shoshtari/paroo/internal/pkg"
)

type MarketRepo interface {
	Insert(context.Context, pkg.Market) (int, error)
}

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
		`
	_, err := m.pool.Exec(ctx, stmt)
	return err

}

func (m MarketsRepoImp) Insert(ctx context.Context, market pkg.Market) (int, error) {
	stmt := `
		INSERT INTO markets(
			exchange_name,
			base_asset,
			quote_asset
		) RETURNING id ON CONFLICT DO NOTHING
		`
	var marketID int
	err := m.pool.QueryRow(ctx, stmt,
		market.ExchangeName,
		market.BaseAsset,
		market.QuoteAsset,
	).Scan(&marketID)

	return marketID, err

}

func NewMarketRepo(pool *pgxpool.Pool, ctx context.Context) (MarketRepo, error) {
	ans := MarketsRepoImp{
		pool: pool,
	}
	return ans, ans.migrate(ctx)
}
