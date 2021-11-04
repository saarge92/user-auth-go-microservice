package containers

import (
	"go-user-microservice/internal/pkg/providers"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
	"go.uber.org/dig"
)

func ProvideShareRepositories(container *dig.Container) error {
	e := container.Provide(
		func(connProvider *providers.ConnectionProvider) *sharedRepositories.CurrencyRepository {
			return sharedRepositories.NewCurrencyRepository(connProvider.GetCoreConnection())
		},
	)
	if e != nil {
		return e
	}
	e = container.Provide(
		func(connProvider *providers.ConnectionProvider) *sharedRepositories.CountryRepository {
			return sharedRepositories.NewCountryRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	return nil
}
