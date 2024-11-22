package wallex

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories/sqlite"
	"github.com/shoshtari/paroo/test"
	"github.com/stretchr/testify/assert"
)

var wallexClient exchange.Exchange
var config configs.ParooConfig

func TestMain(m *testing.M) {
	var err error
	config = test.GetTestConfig()

	marketRepo, err := sqlite.NewMarketRepo(":memory:")
	if err != nil {
		panic(err)
	}

	wallexClient, err = NewWallexClient(config.Wallex, marketRepo)
	if err != nil {
		panic(err)
	}

	if err := pkg.InitializeLogger(config.Log); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestBalance(t *testing.T) {
	balance, err := wallexClient.GetTotalBalance()
	assert.Nil(t, err)
	assert.False(t, balance.Equal(decimal.Zero))
}

func TestMarkets(t *testing.T) {
	stats, err := wallexClient.GetMarkets()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}
func TestMarketStat(t *testing.T) {
	stats, err := wallexClient.GetMarketsStats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}
