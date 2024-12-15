package repositories

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
)

type MarketRepo interface {
	GetOrCreate(context.Context, pkg.Market) (int, bool, error)
	GetByID(context.Context, int) (pkg.Market, error)
	GetByExchangeAndAsset(ctx context.Context, exchange, baseAsset, quoteAsset string) (pkg.Market, error)
	GetAllExchangeMarkets(ctx context.Context, exchangeName string) ([]pkg.Market, error)
}

type BalanceRepo interface {
	Insert(ctx context.Context, changeName string, date time.Time, balance decimal.Decimal) error
	Get(ctx context.Context, exchangeName string, start, date time.Time) ([]time.Time, []decimal.Decimal, error)
}

type MarketStatsRepo interface {
	Insert(ctx context.Context, stat pkg.MarketStat) error
	GetMarketLastStat(ctx context.Context, marketID int) (pkg.MarketStat, error)
}

type ExchangeRepo interface {
	Insert(ctx context.Context, exchange pkg.Exchange) (int, error)
	GetByName(context.Context, string) (int, error)
}
