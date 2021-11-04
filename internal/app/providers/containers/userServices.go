package containers

import (
	"go-user-microservice/internal/app/config"
	stripeInterface "go-user-microservice/internal/app/domain/services/stripe"
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
			accountStripeService stripeInterface.AccountStripeServiceInterface,
		) *user.ServiceUser {
			return user.NewUserService(
				userRepo,
				countryRepository,
				userRemoteService,
				accountStripeService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *repositories.UserRepository,
		) *user.JwtService {
			return user.NewJwtService(config, userRepo)
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
