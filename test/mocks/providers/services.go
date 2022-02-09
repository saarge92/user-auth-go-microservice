package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	userServices "go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/providers"
)

type TestServiceProvider struct {
	AccountStripeMock           stripe.AccountStripeService
	CardStripeServiceMock       stripe.CardStripeService
	AccountStripeServiceMock    stripe.AccountStripeService
	CardStripeChargeServiceMock stripe.ChargeService
	*providers.ServiceProvider
}

func (p *TestServiceProvider) Charge() stripe.ChargeService {
	return p.CardStripeChargeServiceMock
}

func (p *TestServiceProvider) Account() stripe.AccountStripeService {
	return p.AccountStripeMock
}

func (p *TestServiceProvider) Card() stripe.CardStripeService {
	return p.CardStripeServiceMock
}

func (p *TestServiceProvider) AuthService() *userServices.AuthService {
	return p.ServiceProvider.AuthService()
}

func (p *TestServiceProvider) JwtService() *userServices.JwtService {
	return p.ServiceProvider.JwtService()
}

func (p *TestServiceProvider) RemoteUserService() services.RemoteUserService {
	return p.ServiceProvider.RemoteUserService()
}

func (p *TestServiceProvider) UserService() *userServices.ServiceUser {
	return p.ServiceProvider.UserService()
}

func (p *TestServiceProvider) WalletService() services.WalletService {
	return p.ServiceProvider.WalletService()
}

func (p *TestServiceProvider) StripeAccountService() stripe.AccountStripeService {
	if p.AccountStripeServiceMock != nil {
		return p.AccountStripeServiceMock
	}
	return p.ServiceProvider.StripeAccountService()
}

func (p *TestServiceProvider) StripeCardService() stripe.CardStripeService {
	if p.CardStripeServiceMock != nil {
		return p.CardStripeServiceMock
	}
	return p.ServiceProvider.StripeCardService()
}

func (p *TestServiceProvider) UserAuthContextService() *userServices.UserAuthContextService {
	return p.ServiceProvider.UserAuthContextService()
}

func (p *TestServiceProvider) CardService() *cardServices.ServiceCard {
	return p.ServiceProvider.CardService()
}
