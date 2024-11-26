package wallex

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
	"go.uber.org/zap"
)

type ListMarketStatsResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Symbols map[string]struct {
			BaseAsset  string `json:"baseAsset"`
			QuoteAsset string `json:"quoteAsset"`
			Stats      struct {
				BuyPrice  string `json:"bidPrice"`
				SellPrice string `json:"askPrice"`
			} `json:"stats"`
		} `json:"symbols"`
	} `json:"result"`
}

func (w wallexClientImp) GetMarkets() ([]pkg.Market, error) {
	return w.marketsRepo.GetAllExchangeMarkets(context.Background(), exchangeName)
}

func (w wallexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	var res ListMarketStatsResponse
	if err := w.sendReq("markets", nil, &res); err != nil {
		return nil, errors.Wrap(err, "couldn't send request")
	}
	if !res.Success {
		return nil, pkg.InternalError
	}

	var stats []pkg.MarketStat
	for _, symbol := range res.Result.Symbols {
		if _, exists := avoidingSymbols[symbol.BaseAsset]; exists {
			continue
		}
		logger := pkg.GetLogger().With(
			zap.String("exchange", "wallex"),
			zap.String("method", "GetMarketStats"),
			zap.String("base asset", symbol.BaseAsset),
			zap.String("quote asset", symbol.QuoteAsset),
		)

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

		buyprice, err := decimal.NewFromString(symbol.Stats.BuyPrice)
		if err != nil {
			logger.With(zap.Error(err)).Error("couldn't convert buy stats to decimal")
			continue
		}
		sellprice, err := decimal.NewFromString(symbol.Stats.SellPrice)
		if err != nil {
			logger.With(zap.Error(err)).Error("couldn't convert sell stats to decimal")
			continue
		}
		marketStat := pkg.MarketStat{
			MarketID:  market.ID,
			BuyPrice:  buyprice,
			SellPrice: sellprice,
			Date:      time.Now(),
		}
		stats = append(stats, marketStat)
	}

	return stats, nil
}
