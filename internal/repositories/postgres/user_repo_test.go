package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	"github.com/shoshtari/paroo/test"
	"github.com/shoshtari/paroo/test/testcontainer"
	"github.com/stretchr/testify/assert"
)

var config configs.ParooConfig
var userRepo repositories.UserRepo

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

	pool, err := ConnectPostgres(ctx, config.Database.Postgres)
	if err != nil {
		panic(err)
	}

	userRepo, err = NewUserRepo(context.TODO(), pool)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func assertUserEqual(t *testing.T, user1, user2 pkg.User) {
	assert.Equal(t, user1.TelegramID, user2.TelegramID)
	assert.Equal(t, user1.TelegramUsername, user2.TelegramUsername)
	assert.Equal(t, user1.WallexToken, user2.WallexToken)
	assert.Equal(t, user1.RamzinexToken, user2.RamzinexToken)
}
func TestUserRepo(t *testing.T) {
	user := pkg.User{
		TelegramID:       123456,
		TelegramUsername: "john",
		WallexToken:      "wallex_token",
		RamzinexToken:    "ramzinex_token",
	}

	_, err := userRepo.GetOrCreate(context.TODO(), user)
	assert.Nil(t, err)

	user2, err := userRepo.GetOrCreate(context.TODO(), user)
	assert.Nil(t, err)
	assertUserEqual(t, user, user2)
	assert.Nil(t, user2.DeletedAt)
	// created at not timezero
	assert.False(t, user2.CreatedAt.IsZero())

	users, err := userRepo.GetAll(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assertUserEqual(t, user, users[0])
}
