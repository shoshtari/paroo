package wallex

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories/postgres"
	"github.com/shoshtari/paroo/test"
	"github.com/shoshtari/paroo/test/testcontainer"
	"github.com/stretchr/testify/assert"
)

var wallexClient exchange.Exchange
var config configs.ParooConfig
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config = test.GetTestConfig()

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

	wallexClient, err = NewWallexClient(config.Exchange.Wallex, marketRepo)
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

	_, err = pool.Exec(context.Background(), "UPDATE markets SET is_active = TRUE")
	assert.Nil(t, err)

	stats, err := wallexClient.GetMarketsStats()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
	assert.NotEmpty(t, stats)
}
