package providers

import (
	"go-user-microservice/internal/app/card"
	cardDomain "go-user-microservice/internal/app/card/domain"
	paymentDomain "go-user-microservice/internal/app/payment/domain"
	paymentRepositories "go-user-microservice/internal/app/payment/repositories"
	userDomain "go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/repositories"
	"go-user-microservice/internal/app/wallet/domain"
	walletRepositories "go-user-microservice/internal/app/wallet/repositories"
	"go-user-microservice/internal/pkg/domain/providers"
	repositoryInterfaces "go-user-microservice/internal/pkg/domain/repositories"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
)

type RepositoryProvider struct {
	userRepository           userDomain.UserRepository
	currencyRepository       repositoryInterfaces.CurrencyRepository
	walletRepository         domain.WalletRepository
	countryRepository        repositoryInterfaces.CountryRepository
	cardRepository           cardDomain.CardRepository
	operationStoryRepository paymentDomain.OperationStoryRepository
}

func NewRepositoryProvider(dbConnectionProvider providers.DatabaseConnectionProvider) *RepositoryProvider {
	mainDBConnection := dbConnectionProvider.GetCoreConnection()
	return &RepositoryProvider{
		userRepository:           repositories.NewUserRepository(mainDBConnection),
		currencyRepository:       sharedRepositories.NewCurrencyRepository(mainDBConnection),
		walletRepository:         walletRepositories.NewWalletRepository(mainDBConnection),
		countryRepository:        sharedRepositories.NewCountryRepository(mainDBConnection),
		cardRepository:           card.NewRepositoryCard(mainDBConnection),
		operationStoryRepository: paymentRepositories.NewOperationStoryRepository(mainDBConnection),
	}
}

func (p *RepositoryProvider) UserRepository() userDomain.UserRepository {
	return p.userRepository
}

func (p *RepositoryProvider) CurrencyRepository() repositoryInterfaces.CurrencyRepository {
	return p.currencyRepository
}

func (p *RepositoryProvider) WalletRepository() domain.WalletRepository {
	return p.walletRepository
}

func (p *RepositoryProvider) CountryRepository() repositoryInterfaces.CountryRepository {
	return p.countryRepository
}

func (p *RepositoryProvider) CardRepository() cardDomain.CardRepository {
	return p.cardRepository
}

func (p *RepositoryProvider) OperationStory() paymentDomain.OperationStoryRepository {
	return p.operationStoryRepository
}
