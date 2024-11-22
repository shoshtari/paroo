package wallex

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
	"go.uber.org/zap"
)

type BalanceDetailResponse struct {
	Success bool `json:"success"`
	Result  map[string]struct {
		Symbol    string          `json:"symbol"`
		Total     decimal.Decimal `json:"total"`
		Available decimal.Decimal `json:"available"`
	} `json:"result"`
}

func (w wallexClientImp) GetTotalBalance() (decimal.Decimal, error) {
	var wallexRes BalanceDetailResponse
	if err := w.sendReq("account/balances-detail", nil, &wallexRes); err != nil {
		return decimal.Zero, errors.Wrap(err, "couldn't send request")
	}

	if !wallexRes.Success {
		return decimal.Zero, pkg.InternalError
	}

	stats, err := w.GetMarketsStats()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "couldn't get market stats")
	}
	marketPrices := make(map[string]decimal.Decimal)
	marketPrices["TMN"] = decimal.NewFromInt(1)
	for _, stat := range stats {
		market, err := w.marketsRepo.GetByID(context.Background(), stat.MarketID)
		if err != nil {
			return decimal.Zero, errors.Wrap(err, "couldn't get market")
		}
		if market.QuoteAsset != "TMN" {
			continue
		}

		marketPrices[market.BaseAsset] = stat.BuyPrice
	}

	ans := decimal.Zero
	for _, asset := range wallexRes.Result {
		if asset.Total.Equal(decimal.Zero) {
			continue
		}

		price, exists := marketPrices[asset.Symbol]
		if !exists {
			pkg.GetLogger().With(zap.String("symbol", asset.Symbol)).Warn("there was a symbol in portfolio without equivalent in stats")
			continue
		}
		ans = ans.Add(price.Mul(asset.Total))
	}

	return ans, nil
}
