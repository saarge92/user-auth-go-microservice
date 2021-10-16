package containers

import (
	providers2 "go-user-microservice/internal/app/providers"
	repositories2 "go-user-microservice/internal/app/repositories"
	"go.uber.org/dig"
)

func ProvideRepositories(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers2.ConnectionProvider) *repositories2.UserRepository {
			return repositories2.NewUserRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	return nil
}
