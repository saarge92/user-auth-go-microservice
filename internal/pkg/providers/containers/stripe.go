package containers

import (
	"go-user-microservice/internal/pkg/config"
	stripeInterface "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/services"
	"go-user-microservice/internal/pkg/services/stripe"
	"go.uber.org/dig"
)

func ProvideStripeService(container *dig.Container) error {
	e := container.Provide(
		func(
			config *config.Config,
			encryptService *services.EncryptService) *stripe.ClientStripeWrapper {
			return stripe.NewClientStripe(config, encryptService)
		})
	if e != nil {
		return e
	}
	e = container.Provide(
		func(client *stripe.ClientStripeWrapper) stripeInterface.AccountStripeServiceInterface {
			stripeImpl := stripe.NewAccountStripeService(client)
			return stripeImpl
		})
	return e
}
