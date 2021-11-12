package providers

import (
	"go-user-microservice/internal/app/wallet/middlewares"
	"go-user-microservice/internal/pkg/services"
	"go.uber.org/dig"
)

func ProvideGrpcMiddleware(container *dig.Container) error {
	return container.Provide(
		func(authContextService *services.UserAuthContextService) *middlewares.WalletGrpcMiddleware {
			return middlewares.NewWalletGrpcServerMiddleware(authContextService)
		},
	)
}
