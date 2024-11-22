package repositories

import (
	"context"

	"github.com/shoshtari/paroo/internal/pkg"
)

type MarketRepo interface {
	GetOrCreate(context.Context, pkg.Market) (int, error)
	GetByID(context.Context, int) (pkg.Market, error)
	GetAllExchangeMarkets(ctx context.Context, exchangeName string) ([]pkg.Market, error)
}
