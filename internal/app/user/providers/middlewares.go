package providers

import (
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/app/wallet/middlewares"
	"go.uber.org/dig"
)

func ProvideGrpcMiddleware(container *dig.Container) error {
	e := container.Provide(
		func(jwtService *services.JwtService) *middlewares.WalletGrpcMiddleware {
			return middlewares.NewUserGrpcMiddleware(jwtService)
		},
	)
	return e
}
