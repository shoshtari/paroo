package testcontainer

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	ctx       context.Context
	Container *redis.RedisContainer
}

func InitRedis(ctx context.Context) (RedisContainer, error) {
	start := time.Now()

	redisContainer, err := redis.Run(ctx,
		"redis:7",
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("6379/tcp"),
			).WithDeadline(5*time.Second)))

	if err != nil {
		log.WithField("init", "redis").WithError(err).Fatal("failed to start container")
	}
	log.Infof("started redis container in %s", time.Since(start))

	return RedisContainer{
		ctx:       ctx,
		Container: redisContainer,
	}, nil
}

func (p *RedisContainer) Hostname() string {
	host, err := p.Container.Host(p.ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get the host of redis container: %w", err))
	}
	return host
}

func (p *RedisContainer) Port() int {
	mappedPort, err := p.Container.MappedPort(p.ctx, "6379/tcp")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped port of redis container: %w", err))
	}
	return mappedPort.Int()
}

func (p *RedisContainer) Terminate() error {
	return p.Container.Terminate(p.ctx)
}
