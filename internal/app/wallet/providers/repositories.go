package providers

import (
	"go-user-microservice/internal/app/wallet/repositories"
	"go-user-microservice/internal/pkg/providers"
	"go.uber.org/dig"
)

func ProvideWalletRepository(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.WalletRepository {
			return repositories.NewWalletRepository(connProvider.GetCoreConnection())
		},
	)
	return e
}
