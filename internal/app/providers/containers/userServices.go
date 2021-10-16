package containers

import (
	"go-user-microservice/internal/app/config"
	repoInterfaces "go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/repositories"
	"go-user-microservice/internal/app/services/member"
	"go.uber.org/dig"
)

func ProvideUserServices(container *dig.Container) error {
	e := container.Provide(
		func(config *config.Config) *member.RemoteUserService {
			return member.NewRemoteUserService(config)
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userRepo *repositories.UserRepository,
			userRemoteService *member.RemoteUserService,
		) *member.UserService {
			var userRepositoryInterface repoInterfaces.UserRepositoryInterface = userRepo
			return member.NewUserService(userRepositoryInterface, userRemoteService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *repositories.UserRepository,
		) *member.JwtService {
			var userRepositoryInterface repoInterfaces.UserRepositoryInterface = userRepo
			return member.NewJwtService(config, userRepositoryInterface)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userService *member.UserService,
			jwtService *member.JwtService,
		) *member.AuthService {
			return member.NewAuthService(userService, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
