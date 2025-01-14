package exchange

import (
	"github.com/shoshtari/paroo/internal/pkg"
)

type Exchange interface {
	GetPortFolio(user pkg.User) (pkg.PortFolio, error)
	GetMarkets() ([]pkg.Market, error)
	GetMarketsStats() ([]pkg.MarketStat, error)
	GetExchangeInfo() pkg.Exchange
}
