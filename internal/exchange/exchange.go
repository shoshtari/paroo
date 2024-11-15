// now it is only prototype to make my mind more clear
// it is not fixed, or gotten by data, I just implemented what was on my mind
package exchange

import (
	"context"

	"github.com/shopspring/decimal"
)

type OrderType int

const (
	BuyOrder  = iota
	SellOrder = iota
)

func (o OrderType) String() string {
	switch o {
	case 0:
		return "buy"
	case 1:
		return "sell"
	}
	return ""
}

type Market struct {
	Symbol   string
	Exchange Exchange
}

type Order struct {
	Market Market
	Type   OrderType
	Price  decimal.Decimal
	Volume decimal.Decimal
}

type Exchange interface {
	GetMarkets(context.Context) ([]Order, error)
	PlaceOrder(context.Context, Order) error
	RemoveOrder(context.Context, Order) error
}
