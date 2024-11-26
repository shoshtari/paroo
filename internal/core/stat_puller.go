package core

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (p ParooCoreImp) getStatDaemon() {
	var wg errgroup.Group
	for range time.NewTicker(time.Second).C {
		wg.Go(p.getMarketStat)
		wg.Go(p.getWalletStat)
		if err := wg.Wait(); err != nil {
			panic(err)
		}
	}
}

func (p ParooCoreImp) getWalletStat() error {
	portfolio, err := p.wallexClient.GetPortFolio()
	if err != nil {
		return errors.Wrap(err, "couldn't get portfolio from wallex")
	}

	stats, err := p.wallexClient.GetMarketsStats()
	if err != nil {
		return errors.Wrap(err, "couldn't get stats from wallex for wallet")
	}

	markets, err := p.wallexClient.GetMarkets()
	if err != nil {
		return errors.Wrap(err, "couldn't get markets for wallex for wallet")
	}

	marketIDToPrice := make(map[int]decimal.Decimal)
	for _, stat := range stats {
		marketIDToPrice[stat.MarketID] = stat.BuyPrice
	}

	marketSymbolToPrice := make(map[string]decimal.Decimal)
	for _, market := range markets {
		if price, exists := marketIDToPrice[market.ID]; !exists {
			continue
		} else {
			marketSymbolToPrice[market.BaseAsset] = price
		}
	}
	marketSymbolToPrice["TMN"] = decimal.NewFromInt(1)

	balance := decimal.Zero
	for _, asset := range portfolio.Assets {
		price, exists := marketSymbolToPrice[asset.Symbol]
		if !exists {
			pkg.GetLogger().Warn("couldn't find price of asset", zap.String("symbol", asset.Symbol))
			continue
		}

		balance = balance.Add(asset.Value.Mul(price))

	}
	if err := p.balanceRepo.Insert(context.TODO(), "wallex", time.Now(), balance); err != nil {
		pkg.GetLogger().With(zap.String("balance", balance.String())).Error("couldn't insert to db")
		return errors.Wrap(err, "couldn't insert to db")
	}

	return nil
}
func (p ParooCoreImp) getMarketStat() error {
	stats, err := p.wallexClient.GetMarketsStats()
	if err != nil {
		return errors.Wrap(err, "couldn't get market stats from wallex")
	}

	for _, stat := range stats {
		if err := p.statRepo.Insert(context.TODO(), stat); err != nil {
			return errors.Wrap(err, "couldn't insert stat to db")
		}
	}
	return nil
}
