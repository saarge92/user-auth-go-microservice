package containers

import (
	"go-user-microservice/internal/app/providers"
	"go-user-microservice/internal/app/repositories"
	"go.uber.org/dig"
)

func ProvideRepositories(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.UserRepository {
			return repositories.NewUserRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.WalletRepository {
			return repositories.NewWalletRepository(connProvider.GetCoreConnection())
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.CurrencyRepository {
			return repositories.NewCurrencyRepository(connProvider.GetCoreConnection())
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.CountryRepository {
			return repositories.NewCountryRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	return nil
}
