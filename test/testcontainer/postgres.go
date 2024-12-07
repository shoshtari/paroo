package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/shoshtari/paroo/internal/configs"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	ctx       context.Context
	Container *postgres.PostgresContainer
}

func InitPostgres(ctx context.Context, config configs.SectionPostgres) (PostgresContainer, error) {
	start := time.Now()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.User),
		postgres.WithPassword(config.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(5*time.Second)))

	if err != nil {
		log.WithField("init", "postgres").WithError(err).Fatal("failed to start container")
	}
	log.Infof("started postgres container in %s", time.Since(start))

	return PostgresContainer{
		ctx:       ctx,
		Container: postgresContainer,
	}, nil
}

func (p *PostgresContainer) Hostname() string {
	host, err := p.Container.Host(p.ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get the host of postgres container: %w", err))
	}
	return host
}

func (p *PostgresContainer) Port() int {
	mappedPort, err := p.Container.MappedPort(p.ctx, "5432/tcp")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped port of postgres container: %w", err))
	}
	return mappedPort.Int()
}

func (p *PostgresContainer) Terminate() error {
	return p.Container.Terminate(p.ctx)
}
