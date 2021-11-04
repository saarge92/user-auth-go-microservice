package containers

import (
	"go-user-microservice/internal/app/services"
	"go.uber.org/dig"
)

func ProvideEncryptionService(container *dig.Container) error {
	e := container.Provide(
		func() *services.EncryptService {
			return &services.EncryptService{}
		})
	return e
}
