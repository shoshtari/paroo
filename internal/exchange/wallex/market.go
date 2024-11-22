package wallex

import (
	"context"
	"time"

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

func (w wallexClientImp) GetMarkets() ([]pkg.Market, error) {
	return w.marketsRepo.GetAllExchangeMarkets(context.Background(), exchangeName)
}

func (w wallexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	var res ListMarketResponse
	if err := w.sendReq("markets", nil, &res); err != nil {
		return nil, errors.Wrap(err, "couldn't send request")
	}
	if !res.Success {
		return nil, pkg.InternalError
	}

	var stats []pkg.MarketStat
	for _, symbol := range res.Result.Symbols {
		market := pkg.Market{
			ExchangeName: exchangeName,
			BaseAsset:    symbol.BaseAsset,
			QuoteAsset:   symbol.QuoteAsset,
		}
		var err error

		market.ID, err = w.marketsRepo.GetOrCreate(context.TODO(), market)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't insert market to db")
		}

		marketStat := pkg.MarketStat{
			MarketID:  market.ID,
			BuyPrice:  symbol.Stats.BuyPrice,
			SellPrice: symbol.Stats.SellPrice,
			Date:      time.Now(),
		}
		stats = append(stats, marketStat)
	}

	return stats, nil
}
