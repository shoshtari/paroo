package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/shoshtari/paroo/internal/configs"
)

type RedisCache[ValType any] struct {
	client *redis.Client
}

func (r RedisCache[ValType]) encode(val any) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(val); err != nil {
		return nil, errors.Wrap(err, "couldn't encode object")
	}
	return buf.Bytes(), nil
}

func (r RedisCache[ValType]) decode(encodedData []byte, val any) error {
	buf := bytes.NewBuffer(encodedData)
	if err := gob.NewDecoder(buf).Decode(val); err != nil {
		return errors.Wrap(err, "couldn't decode object")
	}
	return nil
}

func (i RedisCache[ValType]) Get(key string) (ValType, error) {
	panic("not implemented")
}

func (i RedisCache[ValType]) Exists(key string) (bool, error) {
	panic("not implemented")
}

func (i RedisCache[ValType]) Set(key string, val ValType) error {
	panic("not implemented")
}

func NewRedisCache[valtype any](config configs.SectionRedis) (Cache[string, valtype], error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.Host, config.Port),
		DB:   config.DB,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, errors.Wrap(err, "couldn't ping redis")
	}
	c := RedisCache[valtype]{
		client: client,
	}
	return c, nil
}
