package ramzinex

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
)

type accountDetailResponseItem struct {
	ID          int    `json:"id"`
	AccountType string `json:"account_type"`
	Name        string `json:"name"`
	Balances    []struct {
		AccountID       int             `json:"account_id"`
		CurrencyID      int             `json:"currency_id"`
		TotalAmount     decimal.Decimal `json:"total_amount"`
		AvailableAmount decimal.Decimal `json:"available_amount"`
		DebtAmount      decimal.Decimal `json:"debt_amount"`
	} `json:"balances"`
}

func (r ramzinexClientImp) GetPortFolio() (pkg.PortFolio, error) {
	ans := pkg.PortFolio{
		ExchangeName: exchangeName,
	}

	var res []accountDetailResponseItem
	err := r.sendReq(sendReqRequest{
		path:    "wallet/api/v1/accounts/detail",
		reqbody: nil,
		resbody: &res,
		auth:    true,
	})
	if err != nil {
		return ans, errors.Wrap(err, "couldn't send request to ramzinex")
	}
	balances := make(map[int]decimal.Decimal)
	for _, item := range res {
		for _, balance := range item.Balances {
			if _, exists := balances[balance.CurrencyID]; !exists {
				balances[balance.CurrencyID] = decimal.Zero
			}
			balances[balance.CurrencyID] = balances[balance.CurrencyID].Add(balance.TotalAmount).Sub(balance.DebtAmount)
		}
	}
	for currencyID, value := range balances {
		currencyName, err := r.currencyIDToName(currencyID)
		if err != nil {
			return ans, errors.Wrap(err, "couldn't convert currency id to name")
		}
		ans.Assets = append(ans.Assets, pkg.Asset{
			Symbol: currencyName,
			Value:  value,
		})

	}
	return ans, nil
}
