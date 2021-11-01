package containers

import (
	"go-user-microservice/internal/app/config"
	repoInterfaces "go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/repositories"
	"go-user-microservice/internal/app/services/user"
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
			countryRepository *repositories.CountryRepository,
		) *user.ServiceUser {
			var userRepositoryInterface repoInterfaces.UserRepositoryInterface = userRepo
			return user.NewUserService(userRepositoryInterface, countryRepository, userRemoteService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *repositories.UserRepository,
		) *user.JwtService {
			var userRepositoryInterface repoInterfaces.UserRepositoryInterface = userRepo
			return user.NewJwtService(config, userRepositoryInterface)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userService *user.ServiceUser,
			jwtService *user.JwtService,
		) *user.AuthService {
			return user.NewAuthService(userService, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
