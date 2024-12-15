package ramzinex

import (
	"fmt"
	"net/http"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
)

type ramzinexClientImp struct {
	httpClient  http.Client
	baseAddress string
	token       string
	marketRepo  repositories.MarketRepo
}

const exchangeName = "ramzinex"

func (r ramzinexClientImp) GetMarketsStats() ([]pkg.MarketStat, error) {
	panic("unimplemented")
}

func (r ramzinexClientImp) GetPortFolio() (pkg.PortFolio, error) {
	panic("unimplemented")
}

func (w ramzinexClientImp) sendReq(path string, reqbody any, resbody any) error {
	url := fmt.Sprintf("%v/%v", w.baseAddress, path)
	return pkg.SendHTTPRequest(w.httpClient, url, reqbody, resbody, pkg.WithHeader("Authorization", w.token))
}

func NewRamzinexClient(config configs.SectionRamzinex, marketRepo repositories.MarketRepo) (exchange.Exchange, error) {
	ans := ramzinexClientImp{
		baseAddress: config.BaseAddress,
		token:       config.Token,
		marketRepo:  marketRepo,
	}
	ans.httpClient.Timeout = config.Timeout
	return ans, nil
}
