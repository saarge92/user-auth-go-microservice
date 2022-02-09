package providers

import (
	"go-user-microservice/internal/pkg/domain/services/stripe"
)

type TestStripeServiceProvider struct {
	AccountStripeServiceMock stripe.AccountStripeService
	CardStripeServiceMock    stripe.CardStripeService
	CardChargeServiceMock    stripe.ChargeService
}

func (p *TestStripeServiceProvider) Charge() stripe.ChargeService {
	return p.CardChargeServiceMock
}

func (p *TestStripeServiceProvider) Account() stripe.AccountStripeService {
	return p.AccountStripeServiceMock
}

func (p *TestStripeServiceProvider) Card() stripe.CardStripeService {
	return p.CardStripeServiceMock
}
