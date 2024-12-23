package exchange

import (
	"github.com/shoshtari/paroo/internal/pkg"
)

type Exchange interface {
	GetPortFolio() (pkg.PortFolio, error)
	GetMarkets() ([]pkg.Market, error)
	GetMarketsStats() ([]pkg.MarketStat, error)
	GetExchangeInfo() pkg.Exchange
}
