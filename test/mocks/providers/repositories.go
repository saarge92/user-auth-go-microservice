package providers

import (
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/user/domain"
	domain2 "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/domain/repositories"
	"go-user-microservice/internal/pkg/providers"
)

type TestingRepositoryProvider struct {
	UserRepositoryMock     domain.UserRepositoryInterface
	CurrencyRepositoryMock repositories.CurrencyRepositoryInterface
	WalletRepositoryMock   domain2.WalletRepositoryInterface
	CountryRepositoryMock  repositories.CountryRepositoryInterface
	CardRepositoryMock     cardDomain.CardRepositoryInterface
	*providers.RepositoryProvider
}

func (p *TestingRepositoryProvider) UserRepository() domain.UserRepositoryInterface {
	if p.UserRepositoryMock != nil {
		return p.UserRepositoryMock
	}
	return p.RepositoryProvider.UserRepository()
}

func (p *TestingRepositoryProvider) CurrencyRepository() repositories.CurrencyRepositoryInterface {
	if p.CurrencyRepositoryMock != nil {
		return p.CurrencyRepositoryMock
	}
	return p.RepositoryProvider.CurrencyRepository()
}

func (p *TestingRepositoryProvider) WalletRepository() domain2.WalletRepositoryInterface {
	if p.WalletRepositoryMock != nil {
		return p.WalletRepositoryMock
	}
	return p.RepositoryProvider.WalletRepository()
}

func (p *TestingRepositoryProvider) CountryRepository() repositories.CountryRepositoryInterface {
	if p.CountryRepositoryMock != nil {
		return p.CountryRepositoryMock
	}
	return p.RepositoryProvider.CountryRepository()
}

func (p *TestingRepositoryProvider) CardRepository() cardDomain.CardRepositoryInterface {
	if p.CardRepositoryMock != nil {
		return p.CardRepositoryMock
	}
	return p.RepositoryProvider.CardRepository()
}
