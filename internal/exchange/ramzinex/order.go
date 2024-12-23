package ramzinex

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
	"go.uber.org/zap"
)

type orderBookItem struct {
	Buys  [][]any `json:"buys"`
	Sells [][]any `json:"sells"`
}

// an order is like price, amount, rial_amount, idk, idk, idk, date

func (r ramzinexClientImp) getOrders() (map[int]orderBookItem, error) {
	res := make(map[int]orderBookItem)
	if err := r.sendReq(sendReqRequest{
		path:         "exchange/api/v1.0/exchange/orderbooks/buys_sells",
		reqbody:      nil,
		resbody:      &res,
		usePublicApi: true,
	}); err != nil {
		return nil, errors.Wrap(err, "couldn't get orders from ramzinex")
	}
	return res, nil

}

func (r ramzinexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	orderBook, err := r.getOrders()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get order books")
	}
	buyprice := make(map[int]decimal.Decimal)
	sellprice := make(map[int]decimal.Decimal)
	for pairID, pairOrders := range orderBook {
		for _, order := range pairOrders.Buys {
			orderPrice := decimal.NewFromFloat(order[0].(float64))
			value, exists := buyprice[pairID]
			if !orderPrice.Equal(decimal.Zero) && (!exists || orderPrice.GreaterThan(value)) {
				buyprice[pairID] = orderPrice
			}
		}
		for _, order := range pairOrders.Sells {
			orderPrice := decimal.NewFromFloat(order[0].(float64))
			value, exists := sellprice[pairID]
			if !orderPrice.Equal(decimal.Zero) && (!exists || orderPrice.LessThan(value)) {
				sellprice[pairID] = orderPrice
			}
		}
	}

	var ans []pkg.MarketStat
	for pairID := range buyprice {
		if _, exists := sellprice[pairID]; !exists {
			continue
		}
		marketID, err := r.getMarketIDFromPairID(pairID)
		if err != nil {
			pkg.GetLogger().Error("couldn't get market id from pair id", zap.Error(err))
			continue
		}

		ans = append(ans, pkg.MarketStat{
			MarketID:  marketID,
			BuyPrice:  buyprice[pairID],
			SellPrice: sellprice[pairID],
			Date:      time.Now(),
		})
	}
	return ans, nil

}
