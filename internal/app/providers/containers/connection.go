package containers

import (
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/providers"
	"go.uber.org/dig"
)

type ConnectionProvider struct{}

func (p *ConnectionProvider) Provide(
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
