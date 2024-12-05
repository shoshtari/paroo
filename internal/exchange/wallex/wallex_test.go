package wallex

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories/sqlite"
	"github.com/shoshtari/paroo/test"
	"github.com/stretchr/testify/assert"
)

var wallexClient exchange.Exchange
var config configs.ParooConfig
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	config = test.GetTestConfig()

	db, err = sqlite.Connect(":memory:")
	if err != nil {
		panic(err)
	}

	marketRepo, err := sqlite.NewMarketRepo(context.TODO(), db)
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
	portfolio, err := wallexClient.GetPortFolio()
	assert.Nil(t, err)
	assert.NotZero(t, len(portfolio.Assets))
	assert.NotZero(t, portfolio.Assets[0].Value)
}

func TestMarkets(t *testing.T) {
	markets, err := wallexClient.GetMarkets()
	assert.Nil(t, err)
	assert.NotEmpty(t, markets)

	_, err = db.Exec("UPDATE markets SET is_active = TRUE")
	assert.Nil(t, err)

	stats, err := wallexClient.GetMarketsStats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.NotEmpty(t, stats)
}
