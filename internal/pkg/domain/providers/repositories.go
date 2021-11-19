package providers

import (
	"go-user-microservice/internal/app/card/domain"
	userDomain "go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/pkg/domain/repositories"
)

type RepositoryProvider interface {
	UserRepository() userDomain.UserRepositoryInterface
	CurrencyRepository() repositories.CurrencyRepositoryInterface
	WalletRepository() repositories.WalletRepositoryInterface
	CountryRepository() repositories.CountryRepositoryInterface
	CardRepository() domain.CardRepositoryInterface
}
