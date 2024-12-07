package core

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	"go.uber.org/zap"
)

type GetPriceRequest struct {
	OrderType     pkg.OrderType
	marketID      int
	TimePercision time.Duration // optional
}

type PriceManager interface {
	GetPrice(context.Context, GetPriceRequest) (decimal.Decimal, error)
}

type PriceManagerImp struct {
	logger           *zap.Logger
	statsRepo        repositories.MarketStatsRepo
	defaultPercision time.Duration
}

func (p PriceManagerImp) GetPrice(ctx context.Context, req GetPriceRequest) (decimal.Decimal, error) {
	if req.TimePercision == 0 {

		req.TimePercision = p.defaultPercision
	}
	stat, err := p.statsRepo.GetMarketLastStat(ctx, req.marketID)
	if err != nil {
		return decimal.Zero, errors.WithStack(err)
	}

	if time.Since(stat.Date) > req.TimePercision {
		errorText := fmt.Sprintf("couldn't satisfy %v percision. the best we've got is %v old", req.TimePercision, time.Since(stat.Date))
		return decimal.Zero, errors.Wrap(pkg.InternalError, errorText)
	}

	switch req.OrderType {
	case pkg.BuyOrder:
		return stat.SellPrice, nil
	case pkg.SellOrder:
		return stat.BuyPrice, nil
	default:
		return decimal.Zero, errors.Wrap(pkg.InternalError, "order type is unknown")
	}
}

func NewPriceManager(statsRepo repositories.MarketStatsRepo, logger *zap.Logger) PriceManager {
	return PriceManagerImp{
		statsRepo: statsRepo,
		logger:    logger,
	}
}
