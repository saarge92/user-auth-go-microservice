package providers

import "go-user-microservice/internal/pkg/domain/services/stripe"

type StripeServiceProvider interface {
	Account() stripe.AccountStripeService
	Card() stripe.CardStripeService
	Charge() stripe.ChargeService
}
