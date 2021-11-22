package providers

import "go-user-microservice/internal/pkg/domain/services/stripe"

type StripeServiceProviderInterface interface {
	Account() stripe.AccountStripeServiceInterface
	Card() stripe.CardStripeServiceInterface
}
