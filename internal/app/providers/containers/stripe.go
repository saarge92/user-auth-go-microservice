package containers

import (
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/services"
	"go-user-microservice/internal/app/services/stripe"
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
		func(client *stripe.ClientStripeWrapper) *stripe.AccountStripeService {
			return stripe.NewAccountStripeService(client)
		})
	return e
}
