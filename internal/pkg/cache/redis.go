package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

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

func (r RedisCache[ValType]) decode(encodedData []byte, val *ValType) error {
	buf := bytes.NewBuffer(encodedData)
	if err := gob.NewDecoder(buf).Decode(val); err != nil {
		return errors.Wrap(err, "couldn't decode object")
	}
	return nil
}

func (r RedisCache[ValType]) Get(key string) (ValType, error) {
	redisRes, err := r.client.Get(context.TODO(), key).Result()
	var ans ValType
	if err != nil {
		return ans, errors.Wrap(err, "couldn't get value from redis")
	}
	if err := r.decode([]byte(redisRes), &ans); err != nil {
		return ans, errors.Wrap(err, "couldn't decode res from redis")
	}
	return ans, nil
}

func (r RedisCache[ValType]) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(context.TODO(), key).Result()
	return exists == 0, err
}

func (r RedisCache[ValType]) Set(key string, val ValType) error {
	encodedVal, err := r.encode(val)
	if err != nil {
		return errors.Wrap(err, "couldn't encode val")
	}

	if err := r.client.Set(context.TODO(), key, encodedVal, time.Duration(0)).Err(); err != nil {
		return errors.Wrap(err, "couldn't set value to redis")
	}
	return nil
}

func (r RedisCache[ValType]) Delete(key string) error {
	if err := r.client.Del(context.TODO(), key).Err(); err != nil {
		return errors.Wrap(err, "couldn't delete key from redis")
	}
	return nil
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
