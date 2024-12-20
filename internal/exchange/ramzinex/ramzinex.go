package ramzinex

import (
	"fmt"
	"net/http"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/pkg/cache"
	"github.com/shoshtari/paroo/internal/repositories"
)

type caches struct {
	currencyIDToName cache.Cache[int, string]
}

type ramzinexClientImp struct {
	httpClient  http.Client
	baseAddress string
	token       string
	marketRepo  repositories.MarketRepo
	caches      caches
}

const exchangeName = "ramzinex"

func (r ramzinexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	panic("unimplemented")
}

func (w ramzinexClientImp) sendReq(path string, reqbody any, resbody any, auth bool) error {
	url := fmt.Sprintf("%v/%v", w.baseAddress, path)
	if auth {
		return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody,
			pkg.WithHeader("Authorization", w.token),
		)
	}
	return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody)
}

func (w ramzinexClientImp) GetExchangeInfo() pkg.Exchange {
	return pkg.Exchange{
		Name:         exchangeName,
		RialSymbol:   "irr",
		TetherSymbol: "usdt",
	}
}

func NewRamzinexClient(config configs.SectionRamzinex, marketRepo repositories.MarketRepo) (exchange.Exchange, error) {
	ans := ramzinexClientImp{
		baseAddress: config.BaseAddress,
		token:       config.Token,
		marketRepo:  marketRepo,
		caches: caches{
			currencyIDToName: cache.NewInmemoryCache[int, string](),
		},
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, nil
}
