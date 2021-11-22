package providers

import (
	"go-user-microservice/internal/pkg/domain/services/stripe"
	stripeServices "go-user-microservice/internal/pkg/services/stripe"
)

type StripeServiceProvider struct {
	accountService stripe.AccountStripeServiceInterface
	cardService    stripe.CardStripeServiceInterface
}

func NewStripeServiceProvider(
	stripeClientProvider *ClientStripeProvider,
) *StripeServiceProvider {
	accountService := stripeServices.NewAccountStripeService(stripeClientProvider.MainClient())
	cardService := stripeServices.NewCardStripeService(stripeClientProvider.MainClient())
	return &StripeServiceProvider{
		accountService: accountService,
		cardService:    cardService,
	}
}

func (p *StripeServiceProvider) Account() stripe.AccountStripeServiceInterface {
	return p.accountService
}

func (p *StripeServiceProvider) Card() stripe.CardStripeServiceInterface {
	return p.cardService
}
