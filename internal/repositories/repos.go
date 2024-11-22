package repositories

import (
	"context"

	"github.com/shoshtari/paroo/internal/pkg"
)

type MarketRepo interface {
	Insert(context.Context, pkg.Market) (int, error)
}
