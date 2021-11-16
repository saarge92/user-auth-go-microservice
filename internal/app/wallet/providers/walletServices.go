package providers

import (
	userRepositories "go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/wallet/repositories"
	"go-user-microservice/internal/app/wallet/services"
	"go-user-microservice/internal/pkg/providers"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
	"go.uber.org/dig"
)

func ProvideWalletServices(container *dig.Container) error {
	e := container.Provide(
		func(
			currencyRepo *sharedRepositories.CurrencyRepository,
			walletRepo *repositories.WalletRepository,
			userRepo *userRepositories.UserRepository,
			connectionProvider *providers.ConnectionProvider,
		) *services.WalletService {
			return services.NewWalletService(walletRepo, userRepo, currencyRepo, connectionProvider.GetCoreConnection())
		},
	)
	return e
}
