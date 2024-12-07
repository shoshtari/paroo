package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type marketStatsRepoImp struct {
	pool *pgxpool.Pool
}

func (m marketStatsRepoImp) migrate(ctx context.Context) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS market_stats(
			market_id SERIAL,
			buy_price TEXT,
			sell_price TEXT,
			date TIMESTAMP,
			PRIMARY KEY(market_id, date)
		)
		`
	_, err := m.pool.Exec(ctx, stmt)
	return err

}

func (m marketStatsRepoImp) Insert(ctx context.Context, stat pkg.MarketStat) error {
	stmt := `
		INSERT INTO market_stats(
			market_id,
			buy_price,
			sell_price,
			date
		) VALUES (
			$1, $2, $3, $4
		)
		`
	_, err := m.pool.Exec(ctx, stmt, stat.MarketID, stat.BuyPrice, stat.SellPrice, stat.Date)

	return err

}

func (m marketStatsRepoImp) GetMarketLastStat(ctx context.Context, marketID int) (pkg.MarketStat, error) {
	stmt := `
		SELECT
			buy_price,
			sell_price,
			date
		FROM market_stats WHERE market_id = $1
			ORDER BY DATE DESC limit 1
		`

	stat := pkg.MarketStat{
		MarketID: marketID,
	}

	err := m.pool.QueryRow(ctx, stmt, marketID).Scan(&stat.BuyPrice, &stat.SellPrice, &stat.Date)
	return stat, err
}

func NewMarketStatsRepo(pool *pgxpool.Pool, ctx context.Context) (repositories.MarketStatsRepo, error) {
	ans := marketStatsRepoImp{
		pool: pool,
	}
	return ans, ans.migrate(ctx)
}
