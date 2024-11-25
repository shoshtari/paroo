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

func (p ParooCoreImp) getStatDaemon() error {
	var wg errgroup.Group
	for range time.NewTicker(time.Second).C {
		wg.Go(p.getMarketStat)
		wg.Go(p.getWalletStat)
		if err := wg.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (p ParooCoreImp) getWalletStat() error {
	portfolio, err := p.wallexClient.GetPortFolio()
	if err != nil {
		return errors.Wrap(err, "couldn't get portfolio from wallex")
	}

	stats, err := p.wallexClient.GetMarketsStats()
	if err != nil {
		return errors.Wrap(err, "couldn't get stats from wallex")
	}
	marketIDToPrice := make(map[int]decimal.Decimal)
	for _, stat := range stats {
		marketIDToPrice[stat.MarketID] = stat.BuyPrice
	}

	markets, err := p.wallexClient.GetMarkets()
	if err != nil {
		return errors.Wrap(err, "couldn't get markets for wallex")
	}

	marketSymbolToID := make(map[string]int)
	for _, market := range markets {
		marketSymbolToID[market.BaseAsset] = market.ID
	}

	balance := decimal.Zero
	for _, asset := range portfolio.Assets {
		marketID, exists := marketSymbolToID[asset.Symbol]
		if !exists {
			pkg.GetLogger().Warn("found an asset without market id", zap.String("symbol", asset.Symbol))
			continue
		}

		price, exists := marketIDToPrice[marketID]
		if !exists {
			pkg.GetLogger().Warn("found an asset without market stat", zap.String("symbol", asset.Symbol))
			continue
		}

		balance = balance.Add(asset.Value.Mul(price))

	}
	if err := p.balanceRepo.Insert(context.TODO(), "wallex", time.Now(), balance); err != nil {
		return errors.Wrap(err, "couldn't insert to db")
	}

	return nil
}
func (p ParooCoreImp) getMarketStat() error {
	panic("unimplemented")
}
