package providers

import (
	"go-user-microservice/internal/app/card"
	cardDomain "go-user-microservice/internal/app/card/domain"
	userDomain "go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/wallet/domain"
	walletRepositories "go-user-microservice/internal/app/wallet/repositories"
	"go-user-microservice/internal/pkg/domain/providers"
	repositoryInterfaces "go-user-microservice/internal/pkg/domain/repositories"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
)

type RepositoryProvider struct {
	userRepository     userDomain.UserRepositoryInterface
	currencyRepository repositoryInterfaces.CurrencyRepositoryInterface
	walletRepository   domain.WalletRepositoryInterface
	countryRepository  repositoryInterfaces.CountryRepositoryInterface
	cardRepository     cardDomain.CardRepositoryInterface
}

func NewRepositoryProvider(dbConnectionProvider providers.DatabaseConnectionProviderInterface) *RepositoryProvider {
	mainDBConnectionProvider := dbConnectionProvider.GetCoreConnection()
	return &RepositoryProvider{
		userRepository:     repositories.NewUserRepository(mainDBConnectionProvider),
		currencyRepository: sharedRepositories.NewCurrencyRepository(mainDBConnectionProvider),
		walletRepository:   walletRepositories.NewWalletRepository(mainDBConnectionProvider),
		countryRepository:  sharedRepositories.NewCountryRepository(mainDBConnectionProvider),
		cardRepository:     card.NewRepositoryCard(mainDBConnectionProvider),
	}
}

func (p *RepositoryProvider) UserRepository() userDomain.UserRepositoryInterface {
	return p.userRepository
}

func (p *RepositoryProvider) CurrencyRepository() repositoryInterfaces.CurrencyRepositoryInterface {
	return p.currencyRepository
}

func (p *RepositoryProvider) WalletRepository() domain.WalletRepositoryInterface {
	return p.walletRepository
}

func (p *RepositoryProvider) CountryRepository() repositoryInterfaces.CountryRepositoryInterface {
	return p.countryRepository
}

func (p *RepositoryProvider) CardRepository() cardDomain.CardRepositoryInterface {
	return p.cardRepository
}
