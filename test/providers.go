package test

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/providers"
	"go.uber.org/dig"
	"log"
)

func ProvideTestConnections(
	container *dig.Container,
	driverDB string,
) error {
	e := container.Provide(func(config *config.Config) *providers.ConnectionProvider {
		return NewConnectionTestProvider(config, driverDB)
	})
	if e != nil {
		return e
	}
	return nil
}

func NewConnectionTestProvider(
	config *config.Config,
	driverDB string,
) *providers.ConnectionProvider {
	coreConn, e := sqlx.Open(driverDB, config.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	return providers.NewConnectionProviderForConnection(coreConn)
}
