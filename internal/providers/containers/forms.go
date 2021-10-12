package containers

import (
	"go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/forms/builders"
	repositoriesImpl "go-user-microservice/internal/repositories"
	"go.uber.org/dig"
)

func ProvideForms(container *dig.Container) error {
	e := container.Provide(func(userRepository *repositoriesImpl.UserRepository) *builders.UserFormBuilder {
		var userRepoInterface repositories.UserRepositoryInterface = userRepository
		return builders.NewUserFormBuilder(userRepoInterface)
	})
	if e != nil {
		return e
	}
	return nil
}
