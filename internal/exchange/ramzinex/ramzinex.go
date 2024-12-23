package ramzinex

import (
	"net/http"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/pkg/cache"
	"github.com/shoshtari/paroo/internal/repositories"
)

type caches struct {
	currencyIDToName cache.Cache[int, string]
	pairIDToMarketID cache.Cache[int, int]
}

type ramzinexClientImp struct {
	httpClient        http.Client
	baseAddress       string
	basePublicAddress string
	token             string
	marketRepo        repositories.MarketRepo
	caches            caches
}

const exchangeName = "ramzinex"

func (w ramzinexClientImp) GetExchangeInfo() pkg.Exchange {
	return pkg.Exchange{
		Name:         exchangeName,
		RialSymbol:   "irr",
		TetherSymbol: "usdt",
	}
}

func NewRamzinexClient(config configs.SectionRamzinex, marketRepo repositories.MarketRepo) (exchange.Exchange, error) {
	ans := ramzinexClientImp{
		baseAddress:       config.BaseAddress,
		basePublicAddress: config.BasePublicAddress,
		token:             config.Token,
		marketRepo:        marketRepo,
		caches: caches{
			currencyIDToName: cache.NewInmemoryCache[int, string](),
			pairIDToMarketID: cache.NewInmemoryCache[int, int](),
		},
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, nil
}
