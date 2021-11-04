package providers

import (
	repositories2 "go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/wallet/repositories"
	"go-user-microservice/internal/app/wallet/services"
	repositories3 "go-user-microservice/internal/pkg/repositories"
	"go.uber.org/dig"
)

func ProvideWalletServices(container *dig.Container) error {
	e := container.Provide(
		func(
			currencyRepo *repositories3.CurrencyRepository,
			walletRepo *repositories.WalletRepository,
			userRepo *repositories2.UserRepository,
		) *services.WalletService {
			return services.NewWalletService(walletRepo, userRepo, currencyRepo)
		},
	)
	return e
}
