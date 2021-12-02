package providers

import (
	"go-user-microservice/internal/pkg/domain/services/stripe"
	stripeServices "go-user-microservice/internal/pkg/services/stripe"
)

type StripeServiceProvider struct {
	accountService stripe.AccountStripeService
	cardService    stripe.CardStripeService
	chargeService  stripe.ChargeService
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

func (p *StripeServiceProvider) Account() stripe.AccountStripeService {
	return p.accountService
}

func (p *StripeServiceProvider) Card() stripe.CardStripeService {
	return p.cardService
}

func (p *StripeServiceProvider) Charge() stripe.ChargeService {
	return p.chargeService
}
