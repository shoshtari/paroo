package wallex

import (
	"os"
	"testing"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/repositories/sqlite"
	"github.com/shoshtari/paroo/test"
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
	os.Exit(m.Run())
}

// func TestBalance(t *testing.T) {
// 	balance, err := wallexClient.GetTotalBalance()
// 	assert.Nil(t, err)
// 	assert.False(t, balance.Equal(decimal.Zero))
// }

// func TestMarketStat(t *testing.T) {
// 	stats, err := wallexClient.GetMarketsStats()
// 	assert.Nil(t, err)
// 	assert.NotNil(t, stats)
// }
