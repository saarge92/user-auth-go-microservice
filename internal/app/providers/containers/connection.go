package containers

import (
	config2 "go-user-microservice/internal/app/config"
	providers2 "go-user-microservice/internal/app/providers"
	"go.uber.org/dig"
)

func ProvideConnections(
	container *dig.Container,
) error {
	e := container.Provide(func(config *config2.Config) *providers2.ConnectionProvider {
		return providers2.NewConnectionProvider(config)
	})
	if e != nil {
		return e
	}
	return nil
}
