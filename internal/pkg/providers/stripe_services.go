package providers

import (
	stripeServices "go-user-microservice/internal/pkg/services/stripe"
)

type StripeServiceProvider struct {
	accountService *stripeServices.AccountStripeService
	cardService    *stripeServices.CardStripeService
	chargeService  *stripeServices.ChargeStripeService
}

func NewStripeServiceProvider(
	stripeClientProvider *ClientStripeProvider,
) *StripeServiceProvider {
	accountService := stripeServices.NewAccountStripeService(stripeClientProvider.MainClient())
	cardService := stripeServices.NewCardStripeService(stripeClientProvider.MainClient())
	chargeService := stripeServices.NewChargeStripeService(stripeClientProvider.MainClient())
	return &StripeServiceProvider{
		accountService: accountService,
		cardService:    cardService,
		chargeService:  chargeService,
	}
}

func (p *StripeServiceProvider) Account() *stripeServices.AccountStripeService {
	return p.accountService
}

func (p *StripeServiceProvider) Card() *stripeServices.CardStripeService {
	return p.cardService
}

func (p *StripeServiceProvider) Charge() *stripeServices.ChargeStripeService {
	return p.chargeService
}
