package containers

import (
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/providers"
	"go.uber.org/dig"
)

func ProvideConnections(
	container *dig.Container,
	driverDB string,
) error {
	e := container.Provide(func(config *config.Config) *providers.ConnectionProvider {
		return providers.NewConnectionProvider(config, driverDB)
	})
	if e != nil {
		return e
	}
	return nil
}
