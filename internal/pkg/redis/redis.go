package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/shoshtari/paroo/internal/configs"
)

type RedisHandler interface {
}
type RedisHandlerImp struct {
	client *redis.Client
}

func NewRedisHandler(ctx context.Context, config configs.SectionRedis) RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.Host, config.Port),
		DB:   config.DB,
	})

	return RedisHandlerImp{
		client: client,
	}

}
