package providers

import (
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/app/wallet/services"
	"go.uber.org/dig"
)

func ProvideWalletGrpcServer(container *dig.Container) error {
	return container.Provide(
		func(
			walletService *services.WalletService,
		) *wallet.GrpcWalletServer {
			return wallet.NewWalletGrpcServer(walletService)
		},
	)
}
