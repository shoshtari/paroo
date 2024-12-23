package ramzinex

import (
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/pkg/cache"
	"go.uber.org/zap"
)

type getCurrencyRes struct {
	Status int `json:"status"`
	Data   struct {
		Currencies []struct {
			ID     int    `json:"id"`
			Symbol string `json:"symbol"`
		} `json:"currencies"`
	} `json:"data"`
}

func (r ramzinexClientImp) currencyIDToName(currencyID int) (string, error) {
	logger := pkg.GetLogger().With(
		zap.String("module", "ramzinx"),
		zap.String("method", "currency id to name"),
	)

	val, err := r.caches.currencyIDToName.Get(currencyID)
	if err == nil {
		return val, nil
	}
	if err != cache.NotfoundError {
		logger.With(zap.Error(err)).Error("couldn't get data from cache")
	}

	var res getCurrencyRes
	if err := r.sendReq(
		sendReqRequest{
			path:    "exchange/api/v2.0/exchange/currencies",
			reqbody: nil,
			resbody: &res,
		}); err != nil {
		logger.With(zap.Error(err)).Error("couldn't get data from ramzinex")
	}
	if res.Status != 0 {
		logger.Error("response status is not zero")
		return "", pkg.InternalError
	}

	var ans string
	for _, currency := range res.Data.Currencies {
		if err := r.caches.currencyIDToName.Set(currency.ID, currency.Symbol); err != nil {
			logger.With(
				zap.Error(err),
				zap.String("symbol", currency.Symbol),
				zap.Int("currency id", currency.ID),
			).Error("couldn't set currency in cache")
		}
		if currency.ID == currencyID {
			ans = currency.Symbol
		}
	}
	if ans == "" {
		logger.With(zap.Int("currency id", currencyID)).Error("currency not exist in cache and api res")
		return "", pkg.InternalError
	}
	return ans, nil
}
