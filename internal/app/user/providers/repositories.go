package providers

import (
	repositories2 "go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/pkg/providers"
	"go.uber.org/dig"
)

func ProviderUserRepository(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories2.UserRepository {
			return repositories2.NewUserRepository(connProvider.GetCoreConnection())
		})
	return e
}
