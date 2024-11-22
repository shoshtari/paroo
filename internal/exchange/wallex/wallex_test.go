package wallex

import (
	"os"
	"testing"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/test"
)

var wallexClient exchange.Exchange
var config configs.ParooConfig

func TestMain(m *testing.M) {
	var err error
	config = test.GetTestConfig()
	wallexClient, err = NewWallexClient(config.Wallex)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// func TestBalance(t *testing.T) {
// 	balance, err := wallexClient.GetTotalBalance()
// 	assert.NotNil(t, err)
// 	assert.NotZero(t, balance)
// }
