package containers

import (
	"go-user-microservice/internal/app/repositories"
	"go-user-microservice/internal/app/services"
	"go.uber.org/dig"
)

func ProvideWalletServices(container *dig.Container) error {
	e := container.Provide(
		func(
			currencyRepo *repositories.CurrencyRepository,
			walletRepo *repositories.WalletRepository,
			userRepo *repositories.UserRepository,
		) *services.WalletService {
			return services.NewWalletService(walletRepo, userRepo, currencyRepo)
		},
	)
	return e
}
