package containers

import (
	"go-user-microservice/internal/config"
	repositoriesInterface "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/repositories"
	"go-user-microservice/internal/services"
	"go-user-microservice/internal/services/user"
	"go.uber.org/dig"
)

func ProvideUserServices(container *dig.Container) error {
	e := container.Provide(
		func(config *config.Config) *user.RemoteUserService {
			return user.NewRemoteUserService(config)
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userRepo *repositories.UserRepository,
			userRemoteService *user.RemoteUserService,
		) *services.UserService {
			var userRepositoryInterface repositoriesInterface.UserRepositoryInterface = userRepo
			return services.NewUserService(userRepositoryInterface, userRemoteService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *repositories.UserRepository,
		) *user.JwtService {
			var userRepositoryInterface repositoriesInterface.UserRepositoryInterface = userRepo
			return user.NewJwtService(config, userRepositoryInterface)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userService *services.UserService,
			jwtService *user.JwtService,
		) *user.AuthService {
			return user.NewAuthService(userService, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
