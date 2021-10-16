package containers

import (
	config2 "go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/domain/repositories"
	repositories2 "go-user-microservice/internal/app/repositories"
	user2 "go-user-microservice/internal/app/services/member"
	"go.uber.org/dig"
)

func ProvideUserServices(container *dig.Container) error {
	e := container.Provide(
		func(config *config2.Config) *user2.RemoteUserService {
			return user2.NewRemoteUserService(config)
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userRepo *repositories2.UserRepository,
			userRemoteService *user2.RemoteUserService,
		) *user2.UserService {
			var userRepositoryInterface repositories.UserRepositoryInterface = userRepo
			return user2.NewUserService(userRepositoryInterface, userRemoteService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config2.Config,
			userRepo *repositories2.UserRepository,
		) *user2.JwtService {
			var userRepositoryInterface repositories.UserRepositoryInterface = userRepo
			return user2.NewJwtService(config, userRepositoryInterface)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userService *user2.UserService,
			jwtService *user2.JwtService,
		) *user2.AuthService {
			return user2.NewAuthService(userService, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
