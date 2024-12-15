package core

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/exchange/wallex"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	"github.com/shoshtari/paroo/internal/repositories/postgres"
	"github.com/shoshtari/paroo/test"
	"github.com/shoshtari/paroo/test/testcontainer"
	"github.com/stretchr/testify/assert"
)

var pool *pgxpool.Pool

var wallexClient exchange.Exchange
var marketStatsRepo repositories.MarketStatsRepo
var priceManager PriceManager

func TestMain(m *testing.M) {
	var err error
	config := test.GetTestConfig()

	ctx := context.TODO()
	if pgcontainer, err := testcontainer.InitPostgres(ctx, config.Database.Postgres); err != nil {
		panic(err)
	} else {
		config.Database.Postgres.Host = pgcontainer.Hostname()
		config.Database.Postgres.Port = uint16(pgcontainer.Port())
		defer func() {
			if err := pgcontainer.Terminate(); err != nil {
				panic(err)
			}
		}()
	}

	pool, err = postgres.ConnectPostgres(ctx, config.Database.Postgres)
	if err != nil {
		panic(err)
	}

	marketRepo, err := postgres.NewMarketRepo(context.TODO(), pool)
	if err != nil {
		panic(err)
	}

	marketStatsRepo, err = postgres.NewMarketStatsRepo(context.TODO(), pool)
	if err != nil {
		panic(err)
	}

	wallexClient, err = wallex.NewWallexClient(config.Exchange.Wallex, marketRepo)
	if err != nil {
		panic(err)
	}

	if err := pkg.InitializeLogger(config.Log); err != nil {
		panic(err)
	}
	priceManager = NewPriceManager(marketStatsRepo, pkg.GetLogger("test"))
	os.Exit(m.Run())
}

func TestGetPrice(t *testing.T) {
	markets, err := wallexClient.GetMarkets()
	assert.Nil(t, err)
	assert.NotEmpty(t, markets)

	_, err = priceManager.GetPrice(context.TODO(), GetPriceRequest{
		MarketID:  markets[0].ID,
		OrderType: pkg.SellOrder,
	})
	assert.NotNil(t, err)

	stats, err := wallexClient.GetMarketsStats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.NotEmpty(t, stats)
	for _, stat := range stats {
		err = marketStatsRepo.Insert(context.TODO(), stat)
		assert.Nil(t, err)
	}

	price, err := priceManager.GetPrice(context.TODO(), GetPriceRequest{
		MarketID:  markets[0].ID,
		OrderType: pkg.SellOrder,
	})
	assert.Nil(t, err)
	assert.False(t, price.Equal(decimal.Zero))
}
