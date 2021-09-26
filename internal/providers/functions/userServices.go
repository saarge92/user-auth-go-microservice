package functions

import (
	repositoriesInterface "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/repositories"
	"go-user-microservice/internal/services"
	"go.uber.org/dig"
)

func ProvideUserRepositories(container *dig.Container) error {
	e := container.Provide(func(userRepo *repositories.UserRepository) *services.UserService {
		var userRepositoryInterface repositoriesInterface.UserRepository = userRepo
		return services.NewUserService(userRepositoryInterface)
	})
	if e != nil {
		return e
	}
	return nil
}
