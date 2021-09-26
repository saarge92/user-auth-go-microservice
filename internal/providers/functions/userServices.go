package functions

import (
	"go-user-microservice/internal/config"
	repositoriesInterface "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/repositories"
	"go-user-microservice/internal/services"
	"go.uber.org/dig"
)

func ProvideUserServices(container *dig.Container) error {
	e := container.Provide(
		func(config *config.Config) *services.RemoteUserService {
			return services.NewRemoteUserService(config)
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userRepo *repositories.UserRepository,
			userRemoteService *services.RemoteUserService,
		) *services.UserService {
			var userRepositoryInterface repositoriesInterface.UserRepository = userRepo
			return services.NewUserService(userRepositoryInterface, userRemoteService)
		})
	if e != nil {
		return e
	}
	return nil
}
