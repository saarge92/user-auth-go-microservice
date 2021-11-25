package providers

import (
	"go-user-microservice/internal/app/card/domain"
	userDomain "go-user-microservice/internal/app/user/domain"
	domain2 "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/domain/repositories"
)

type RepositoryProviderInterface interface {
	UserRepository() userDomain.UserRepositoryInterface
	CurrencyRepository() repositories.CurrencyRepositoryInterface
	WalletRepository() domain2.WalletRepositoryInterface
	CountryRepository() repositories.CountryRepositoryInterface
	CardRepository() domain.CardRepositoryInterface
}
