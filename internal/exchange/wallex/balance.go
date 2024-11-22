package wallex

import (
	"fmt"

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

func (w WallexClientImp) GetTotalBalance() (decimal.Decimal, error) {
	var wallexRes BalanceDetailResponse
	if err := w.sendReq("account/balances-detail", nil, &wallexRes); err != nil {
		return decimal.Zero, errors.Wrap(err, "couldn't send request")
	}

	if !wallexRes.Success {
		return decimal.Zero, pkg.InternalError
	}

	for _, asset := range wallexRes.Result {
		if asset.Total.Equal(decimal.Zero) {
			continue
		}
		fmt.Println(asset.Symbol, asset.Total)
	}

	return decimal.Zero, pkg.BadRequestError

}
