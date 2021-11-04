package providers

import (
	repositories2 "go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/repositories"
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
			userRepo *repositories2.UserRepository,
			userRemoteService *services.RemoteUserService,
			countryRepository *repositories.CountryRepository,
		) *services.ServiceUser {
			return services.NewUserService(
				userRepo,
				countryRepository,
				userRemoteService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *repositories2.UserRepository,
		) *services.JwtService {
			return services.NewJwtService(config, userRepo)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			userService *services.ServiceUser,
			jwtService *services.JwtService,
		) *services.AuthService {
			return services.NewAuthService(userService, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
