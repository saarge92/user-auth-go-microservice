package providers

import (
	userRepositories "go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/repositories"
	sharedServices "go-user-microservice/internal/pkg/services"
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
			userRepo *userRepositories.UserRepository,
			userRemoteService *services.RemoteUserService,
			countryRepository *repositories.CountryRepository,
			stripeService stripe.AccountStripeServiceInterface,
		) *services.ServiceUser {
			return services.NewUserService(
				userRepo,
				countryRepository,
				userRemoteService,
				stripeService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			config *config.Config,
			userRepo *userRepositories.UserRepository,
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
	e = container.Provide(
		func(
			jwtService *services.JwtService,
		) *sharedServices.UserAuthContextService {
			return sharedServices.NewUserAuthContextService(jwtService)
		},
	)
	return e
}
