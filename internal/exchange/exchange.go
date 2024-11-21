package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderType int

const (
	BuyOrder OrderType = iota
	SellOrder
)

func (o OrderType) String() string {
	switch o {
	case BuyOrder:
		return "Buy"
	case SellOrder:
		return "Sell"
	default:
		return "Unknown"
	}
}

type Market struct {
	ID           int
	ExchangeName string
	BaseCurrency string
	Currency     string
}

type Order struct {
	Market     Market
	Type       OrderType
	Price      decimal.Decimal
	Amount     decimal.Decimal
	CreateDate time.Time
}

type MarketState struct {
	Market    Market
	SellPrice decimal.Decimal
	BuyPrice  decimal.Decimal

	BuyOrders  []Order
	SellOrders []Order

	Date time.Time
}

type Exchange interface {
	GetTotalBalance() (decimal.Decimal, error)
	GetMarkets() ([]Market, error)
	GetMarketState(market Market) (MarketState, error)
}
