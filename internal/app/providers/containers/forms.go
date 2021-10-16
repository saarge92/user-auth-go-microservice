package containers

import (
	repositories3 "go-user-microservice/internal/app/domain/repositories"
	builders2 "go-user-microservice/internal/app/forms/builders"
	repositories2 "go-user-microservice/internal/app/repositories"
	"go.uber.org/dig"
)

func ProvideForms(container *dig.Container) error {
	e := container.Provide(func(userRepository *repositories2.UserRepository) *builders2.UserFormBuilder {
		var userRepoInterface repositories3.UserRepositoryInterface = userRepository
		return builders2.NewUserFormBuilder(userRepoInterface)
	})
	if e != nil {
		return e
	}
	return nil
}
