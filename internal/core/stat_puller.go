package core

import (
	"context"
	"fmt"
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
			pkg.GetLogger().With(
				zap.String("module", "stat_puller"),
				zap.Error(err),
			).Panic("got error from waitgroup")
		}
	}
}

func (p ParooCoreImp) getWalletStat() error {
	portfolio, err := p.wallexClient.GetPortFolio()
	if err != nil {
		return errors.Wrap(err, "couldn't get portfolio from wallex")
	}

	balance := decimal.Zero
	for _, asset := range portfolio.Assets {
		ctx := context.TODO()
		market, err := p.marketsRepo.GetByExchangeAndAsset(ctx, "wallex", asset.Symbol, "TMN")
		if err != nil {
			return errors.WithStack(err)
		}
		if !market.IsActive {
			continue
		}

		price, err := p.priceManager.GetPrice(context.TODO(), GetPriceRequest{
			OrderType: pkg.SellOrder,
			MarketID:  market.ID,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		balance = balance.Add(asset.Value.Mul(price))

	}
	if err := p.balanceRepo.Insert(context.TODO(), "wallex", time.Now(), balance); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
func (p ParooCoreImp) getMarketStat() error {
	logger := pkg.GetLogger().With(
		zap.String("package", "core"),
		zap.String("module", "stat puller"),
		zap.String("method", "getMarketStat"),
	)

	stats, err := p.wallexClient.GetMarketsStats()
	if err != nil {
		return errors.Wrap(err, "couldn't get market stats from wallex")
	}

	logger.Info(fmt.Sprintf("got %d stat from wallex", len(stats)))

	for _, stat := range stats {
		if err := p.statRepo.Insert(context.TODO(), stat); err != nil {
			return errors.Wrap(err, "couldn't insert stat to db")
		}
	}
	return nil
}
