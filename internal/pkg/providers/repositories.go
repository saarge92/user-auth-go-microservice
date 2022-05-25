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
	roleRepository           *repositories.Role
}

func NewRepositoryProvider(dbConnectionProvider providers.DatabaseConnectionProvider) *RepositoryProvider {
	mainDBWrapper := dbConnectionProvider.GetCoreConnection()

	return &RepositoryProvider{
		userRepository:           repositories.NewUserRepository(mainDBWrapper),
		currencyRepository:       sharedRepositories.NewCurrencyRepository(mainDBWrapper),
		walletRepository:         walletRepositories.NewWalletRepository(mainDBWrapper),
		countryRepository:        sharedRepositories.NewCountryRepository(mainDBWrapper),
		cardRepository:           card.NewRepositoryCard(mainDBWrapper),
		operationStoryRepository: paymentRepositories.NewOperationStoryRepository(mainDBWrapper),
		roleRepository:           repositories.NewRoleRepository(mainDBWrapper),
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

func (p *RepositoryProvider) RoleRepository() *repositories.Role {
	return p.roleRepository
}
