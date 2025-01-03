package wallex

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
)

type BalanceDetailResponse struct {
	Success bool `json:"success"`
	Result  map[string]struct {
		Symbol    string          `json:"symbol"`
		Total     decimal.Decimal `json:"total"`
		Available decimal.Decimal `json:"available"`
	} `json:"result"`
}

func (w wallexClientImp) GetPortFolio() (pkg.PortFolio, error) {
	var wallexRes BalanceDetailResponse
	var ans pkg.PortFolio
	ans.ExchangeName = exchangeName

	if err := w.sendReq("account/balances-detail", nil, &wallexRes, true); err != nil {
		return ans, errors.Wrap(err, "couldn't send request")
	}

	if !wallexRes.Success {
		return ans, pkg.InternalError
	}

	for _, asset := range wallexRes.Result {
		if asset.Total.Equal(decimal.Zero) {
			continue
		}

		ans.Assets = append(ans.Assets, pkg.Asset{
			Symbol: asset.Symbol,
			Value:  asset.Total,
		})
	}

	return ans, nil
}
