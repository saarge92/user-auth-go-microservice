package providers

import (
	"go-user-microservice/internal/pkg/domain/services/stripe"
)

type TestStripeServiceProvider struct {
	AccountStripeServiceMock stripe.AccountStripeServiceInterface
	CardStripeServiceMock    stripe.CardStripeServiceInterface
}

func (p *TestStripeServiceProvider) Account() stripe.AccountStripeServiceInterface {
	return p.AccountStripeServiceMock
}

func (p *TestStripeServiceProvider) Card() stripe.CardStripeServiceInterface {
	return p.CardStripeServiceMock
}
