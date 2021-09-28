package containers

import (
	"go-user-microservice/internal/config"
	"go.uber.org/dig"
)

func ProvideConfig(container *dig.Container) error {
	e := container.Provide(func() *config.Config {
		return config.NewConfig()
	})
	if e != nil {
		return e
	}
	return nil
}
