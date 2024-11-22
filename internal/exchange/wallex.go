package exchange

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/pkg"
)

type WallexClientImp struct {
	baseAddress string
	token       string
	httpClient  http.Client
}

func (w WallexClientImp) sendReq(path string, reqbody any, resbody any) error {
	url := fmt.Sprintf("%v/%v", w.baseAddress, path)
	return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody, pkg.WithHeader("Authorization", w.token))
}

func (w WallexClientImp) getProfile() error {
	return w.sendReq("account/profile", nil, nil)
}

// GetMarketState implements Exchange.
func (WallexClientImp) GetMarketState(market Market) (MarketState, error) {
	panic("unimplemented")
}

// GetMarkets implements Exchange.
func (WallexClientImp) GetMarkets() ([]Market, error) {
	panic("unimplemented")
}

// GetTotalBalance implements Exchange.
func (WallexClientImp) GetTotalBalance() (decimal.Decimal, error) {
	panic("unimplemented")
}

func NewWallexClient(config configs.SectionWallex) (Exchange, error) {
	if config.Token == "" {
		return nil, errors.Wrap(pkg.BadRequestError, "token cannot be empty")
	}
	ans := WallexClientImp{
		baseAddress: config.BaseAddress,
		token:       config.Token,
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, ans.getProfile()
}
