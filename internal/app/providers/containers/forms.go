package containers

import (
	repositories3 "go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/forms/user/builders"
	repositories2 "go-user-microservice/internal/app/repositories"
	"go.uber.org/dig"
)

func ProvideForms(container *dig.Container) error {
	e := container.Provide(func(userRepository *repositories2.UserRepository) *builders.UserFormBuilder {
		var userRepoInterface repositories3.UserRepositoryInterface = userRepository
		return builders.NewUserFormBuilder(userRepoInterface)
	})
	if e != nil {
		return e
	}
	return nil
}
