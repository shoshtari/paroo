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

var avoidingSymbols = map[string]struct{}{
	"BNB":  {},
	"LINK": {},
	"EGLD": {},
}

type wallexClientImp struct {
	baseAddress string
	token       string
	httpClient  http.Client
	marketsRepo repositories.MarketRepo
}

const exchangeName = "wallex"

func (w wallexClientImp) sendReq(path string, reqbody any, resbody any) error {
	url := fmt.Sprintf("%v/%v", w.baseAddress, path)
	return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody, pkg.WithHeader("Authorization", w.token))
}

func (w wallexClientImp) getProfile() error {
	return w.sendReq("account/profile", nil, nil)
}

func NewWallexClient(config configs.SectionWallex, marketsRepo repositories.MarketRepo) (exchange.Exchange, error) {
	if config.Token == "" {
		return nil, errors.Wrap(pkg.BadRequestError, "token cannot be empty")
	}
	ans := wallexClientImp{
		baseAddress: config.BaseAddress,
		token:       config.Token,
		marketsRepo: marketsRepo,
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, ans.getProfile()
}
