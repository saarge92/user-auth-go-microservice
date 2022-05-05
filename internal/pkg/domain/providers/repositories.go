package providers

import (
	"go-user-microservice/internal/app/card/domain"
	paymentDomain "go-user-microservice/internal/app/payment/domain"
	userDomain "go-user-microservice/internal/app/user/domain"
	userRepositories "go-user-microservice/internal/app/user/repositories"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/domain/repositories"
)

type RepositoryProvider interface {
	UserRepository() userDomain.UserRepository
	CurrencyRepository() repositories.CurrencyRepository
	WalletRepository() walletDomain.WalletRepository
	CountryRepository() repositories.CountryRepository
	CardRepository() domain.CardRepository
	OperationStory() paymentDomain.OperationStoryRepository
	RoleRepository() *userRepositories.Role
}
