package pkg

import (
	"time"

	"github.com/shopspring/decimal"
)

type Market struct {
	ID           int
	ExchangeName string
	BaseAsset    string
	QuoteAsset   string
}

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

type Order struct {
	MarketID   int
	Type       OrderType
	Price      decimal.Decimal
	Amount     decimal.Decimal
	CreateDate time.Time
}

type MarketStat struct {
	MarketID  int
	SellPrice decimal.Decimal
	BuyPrice  decimal.Decimal

	Date time.Time
}