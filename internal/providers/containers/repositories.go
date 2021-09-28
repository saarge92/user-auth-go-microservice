package containers

import (
	"go-user-microservice/internal/providers"
	"go-user-microservice/internal/repositories"
	"go.uber.org/dig"
)

func ProvideRepositories(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.UserRepository {
			return repositories.NewUserRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	return nil
}
