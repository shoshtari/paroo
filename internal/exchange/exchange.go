package exchange

import (
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
)

type Exchange interface {
	GetTotalBalance() (decimal.Decimal, error)
	GetMarkets() ([]pkg.Market, error)
	GetMarketsStats() ([]pkg.MarketStat, error)
}
