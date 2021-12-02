package providers

import (
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/user/domain"
	domain2 "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/domain/repositories"
	"go-user-microservice/internal/pkg/providers"
)

type TestingRepositoryProvider struct {
	UserRepositoryMock     domain.UserRepository
	CurrencyRepositoryMock repositories.CurrencyRepository
	WalletRepositoryMock   domain2.WalletRepository
	CountryRepositoryMock  repositories.CountryRepository
	CardRepositoryMock     cardDomain.CardRepository
	*providers.RepositoryProvider
}

func (p *TestingRepositoryProvider) UserRepository() domain.UserRepository {
	if p.UserRepositoryMock != nil {
		return p.UserRepositoryMock
	}
	return p.RepositoryProvider.UserRepository()
}

func (p *TestingRepositoryProvider) CurrencyRepository() repositories.CurrencyRepository {
	if p.CurrencyRepositoryMock != nil {
		return p.CurrencyRepositoryMock
	}
	return p.RepositoryProvider.CurrencyRepository()
}

func (p *TestingRepositoryProvider) WalletRepository() domain2.WalletRepository {
	if p.WalletRepositoryMock != nil {
		return p.WalletRepositoryMock
	}
	return p.RepositoryProvider.WalletRepository()
}

func (p *TestingRepositoryProvider) CountryRepository() repositories.CountryRepository {
	if p.CountryRepositoryMock != nil {
		return p.CountryRepositoryMock
	}
	return p.RepositoryProvider.CountryRepository()
}

func (p *TestingRepositoryProvider) CardRepository() cardDomain.CardRepository {
	if p.CardRepositoryMock != nil {
		return p.CardRepositoryMock
	}
	return p.RepositoryProvider.CardRepository()
}
