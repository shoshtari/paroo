package wallex

import "github.com/shoshtari/paroo/internal/exchange"

type ListMarketResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Symbols map[string]struct {
			Symbol    string `json:"symbol"`
			BaseAsset string `json:"base_asset"`
		} `json:"symbol"`
	} `json:"result"`
}

func (WallexClientImp) GetMarkets() ([]exchange.Market, error) {
	panic("unimplemented")
}

func (WallexClientImp) GetMarketState(market exchange.Market) (exchange.MarketState, error) {
	panic("unimplemented")
}
