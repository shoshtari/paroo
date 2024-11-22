package wallex

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
)

type ListMarketResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Symbols map[string]struct {
			BaseAsset  string `json:"baseAsset"`
			QuoteAsset string `json:"quoteAsset"`
			Stats      struct {
				BuyPrice  decimal.Decimal `json:"bidPrice"`
				SellPrice decimal.Decimal `json:"askPrice"`
			} `json:"stats"`
		} `json:"symbols"`
	} `json:"result"`
}

func (wallexClientImp) GetMarkets() ([]pkg.Market, error) {
	panic("unimplemented")
}

func (w wallexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	var res ListMarketResponse
	if err := w.sendReq("markets", nil, &res); err != nil {
		return nil, errors.Wrap(err, "couldn't send request")
	}
	if !res.Success {
		return nil, pkg.InternalError
	}
	for _, symbol := range res.Result.Symbols {
		fmt.Println(symbol)

	}

	return nil, pkg.NotImplementedError
}
