package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	userServices "go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/providers"
	sharedServices "go-user-microservice/internal/pkg/services"
)

type TestServiceProvider struct {
	AccountStripeMock        stripe.AccountStripeServiceInterface
	CardStripeServiceMock    stripe.CardStripeServiceInterface
	AccountStripeServiceMock stripe.AccountStripeServiceInterface
	*providers.ServiceProvider
}

func (p *TestServiceProvider) Account() stripe.AccountStripeServiceInterface {
	return p.AccountStripeMock
}

func (p *TestServiceProvider) Card() stripe.CardStripeServiceInterface {
	return p.CardStripeServiceMock
}

func (p *TestServiceProvider) AuthService() *userServices.AuthService {
	return p.ServiceProvider.AuthService()
}

func (p *TestServiceProvider) JwtService() *userServices.JwtService {
	return p.ServiceProvider.JwtService()
}

func (p *TestServiceProvider) RemoteUserService() services.RemoteUserServiceInterface {
	return p.ServiceProvider.RemoteUserService()
}

func (p *TestServiceProvider) UserService() *userServices.ServiceUser {
	return p.ServiceProvider.UserService()
}

func (p *TestServiceProvider) WalletService() services.WalletServiceInterface {
	return p.ServiceProvider.WalletService()
}

func (p *TestServiceProvider) StripeAccountService() stripe.AccountStripeServiceInterface {
	if p.AccountStripeServiceMock != nil {
		return p.AccountStripeServiceMock
	}
	return p.ServiceProvider.StripeAccountService()
}

func (p *TestServiceProvider) StripeCardService() stripe.CardStripeServiceInterface {
	if p.CardStripeServiceMock != nil {
		return p.CardStripeServiceMock
	}
	return p.ServiceProvider.StripeCardService()
}

func (p *TestServiceProvider) UserAuthContextService() *sharedServices.UserAuthContextService {
	return p.ServiceProvider.UserAuthContextService()
}

func (p *TestServiceProvider) CardService() *cardServices.ServiceCard {
	return p.ServiceProvider.CardService()
}
