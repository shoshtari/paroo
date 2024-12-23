package wallex

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type wallexClientImp struct {
	baseAddress string
	token       string
	httpClient  http.Client
	marketsRepo repositories.MarketRepo
}

const exchangeName = "wallex"

func (w wallexClientImp) sendReq(path string, reqbody any, resbody any, auth bool) error {
	url := fmt.Sprintf("%v/%v", w.baseAddress, path)
	if auth {
		return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody, pkg.WithHeader("Authorization", w.token))
	}
	return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody)

}

func (w wallexClientImp) getProfile() error {
	return w.sendReq("account/profile", nil, nil, true)
}

func (w wallexClientImp) GetExchangeInfo() pkg.Exchange {
	return pkg.Exchange{
		Name:         exchangeName,
		RialSymbol:   "TMN",
		TetherSymbol: "USDT",
	}
}

func NewWallexClient(config configs.SectionWallex, marketsRepo repositories.MarketRepo) (exchange.Exchange, error) {
	if config.Token == "" {
		return nil, errors.Wrap(pkg.BadRequestError, "token cannot be empty")
	}
	if marketsRepo == nil {
		return nil, errors.Wrap(pkg.BadRequestError, "markets repo is empty")
	}
	ans := wallexClientImp{
		baseAddress: config.BaseAddress,
		token:       config.Token,
		marketsRepo: marketsRepo,
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, ans.getProfile()
}
