package containers

import (
	"go-user-microservice/internal/app/forms/user/builders"
	"go-user-microservice/internal/app/server"
	"go-user-microservice/internal/app/services"
	"go-user-microservice/internal/app/services/user"
	"go.uber.org/dig"
)

func ProvideGrpcServers(container *dig.Container) error {
	e := container.Provide(
		func(
			userService *user.ServiceUser,
			userFormBuilder *builders.UserFormBuilder,
			authService *user.AuthService,
		) *server.UserGrpcServer {
			return server.NewUserGrpcServer(userFormBuilder, authService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(
			walletService *services.WalletService,
		) *server.WalletGrpcServer {
			return server.NewWalletGrpcServer(walletService)
		},
	)
	if e != nil {
		return e
	}
	return nil
}
