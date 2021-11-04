package containers

import (
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/providers"
	"go.uber.org/dig"
)

func ProvideConnection(
	container *dig.Container,
) error {
	e := container.Provide(func(config *config.Config) *providers.ConnectionProvider {
		return providers.NewConnectionProvider(config)
	})
	if e != nil {
		return e
	}
	return nil
}
