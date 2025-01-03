package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/test"
	"github.com/shoshtari/paroo/test/testcontainer"
)

var config configs.ParooConfig

func TestMain(m *testing.M) {

	config = test.GetTestConfig()

	ctx := context.TODO()
	redisContainer, err := testcontainer.InitRedis(ctx)
	if err != nil {
		panic("couldn't init redis container")
	}
	config.Database.Redis.Host = redisContainer.Hostname()
	config.Database.Redis.Port = uint16(redisContainer.Port())
}

func TestRedisCache(t *testing.T) {
	r, err := NewRedisCache[string](config.Database.Redis)
	assert.Nil(t, err)

	err = r.Set("test", "test")
	assert.Nil(t, err)

	val, err := r.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, val, "test")

	err = r.Delete("test")
	assert.Nil(t, err)

	_, err = r.Get("test")
	assert.NotNil(t, err)

}
