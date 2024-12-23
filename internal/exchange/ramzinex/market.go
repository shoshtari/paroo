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

func (r ramzinexClientImp) getPairs() (ramzinexGetPairResponse, error) {
	var ramzinexRes ramzinexGetPairResponse
	err := r.sendReq(sendReqRequest{
		path:    "exchange/api/v2.0/exchange/pairs",
		reqbody: nil,
		resbody: &ramzinexRes,
	})
	if err != nil {
		return ramzinexRes, errors.WithStack(err)
	}

	if ramzinexRes.Status != 0 {
		return ramzinexRes, errors.New("ramzinex status is not zero")
	}

	return ramzinexRes, nil
}

func (r ramzinexClientImp) getMarketIDFromPairID(pairID int) (int, error) {
	if exists, err := r.caches.pairIDToMarketID.Exists(pairID); err != nil {
		return -1, errors.Wrap(err, "couldn't check if pair exists in cache")
	} else if !exists {
		pairs, err := r.getPairs()
		if err != nil {
			return -1, errors.Wrap(err, "couldn't get pairs from ramzinex")
		}
		for _, pair := range pairs.Data.Pairs {
			marketID, _, err := r.marketRepo.GetOrCreate(context.TODO(), pkg.Market{

				ExchangeName: exchangeName,
				BaseAsset:    pair.BaseCurrency.Symbol.En,
				QuoteAsset:   pair.QuoteCurrency.Symbol.En,
			})
			if err != nil {
				return -1, errors.Wrap(err, "couldn't get market id from repo")
			}
			if err := r.caches.pairIDToMarketID.Set(pair.Id, marketID); err != nil {
				return -1, errors.Wrap(err, "couldn't set pair/market in cache")
			}
		}
	}
	if val, err := r.caches.pairIDToMarketID.Get(pairID); err != nil {
		return -1, errors.Wrap(err, "couldn't get market id from cache after filling")
	} else {
		return val, nil
	}
}

func (r ramzinexClientImp) GetMarkets() ([]pkg.Market, error) {
	var ans []pkg.Market
	ramzinexRes, err := r.getPairs()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get pairs from ramzinex")
	}
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
