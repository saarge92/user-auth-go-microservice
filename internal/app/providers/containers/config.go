package containers

import (
	config2 "go-user-microservice/internal/app/config"
	"go.uber.org/dig"
)

func ProvideConfig(container *dig.Container) error {
	e := container.Provide(func() *config2.Config {
		return config2.NewConfig()
	})
	if e != nil {
		return e
	}
	return nil
}
