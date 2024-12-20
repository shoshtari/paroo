package ramzinex

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
)

type ramzinexGetPairResponse struct {
	Status int `json:"int"`
	Data   struct {
		Pairs []struct {
			Id   int `json:"id"`
			Name struct {
				Fa string `json:"fa"`
				En string `json:"en"`
			}
			QuoteCurrency struct {
				Symbol struct {
					En string `json:"en"`
				}
			} `json:"quote_currency"`
			BaseCurrency struct {
				Symbol struct {
					En string `json:"en"`
				}
			} `json:"base_currency"`
		} `json:"pairs"`
	} `json:"data"`
}

func (r ramzinexClientImp) GetMarkets() ([]pkg.Market, error) {
	var ramzinexRes ramzinexGetPairResponse
	err := r.sendReq("exchange/api/v2.0/exchange/pairs", nil, &ramzinexRes, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if ramzinexRes.Status != 0 {
		return nil, errors.New("ramzinex status is not zero")
	}

	var ans []pkg.Market
	for _, pair := range ramzinexRes.Data.Pairs {
		newMarket := pkg.Market{
			ExchangeName: exchangeName,
			BaseAsset:    pair.BaseCurrency.Symbol.En,
			QuoteAsset:   pair.QuoteCurrency.Symbol.En,
			EnName:       pair.Name.En,
			FaName:       pair.Name.Fa,
		}
		newMarket.ID, newMarket.IsActive, err = r.marketRepo.GetOrCreate(context.TODO(), newMarket)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get/create row in db")
		}
		ans = append(ans, newMarket)
	}
	return ans, nil

}
