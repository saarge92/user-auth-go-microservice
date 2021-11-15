package containers

import (
	"go-user-microservice/internal/pkg/config"
	stripeInterface "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/services"
	"go-user-microservice/internal/pkg/services/stripe"
	"go.uber.org/dig"
)

func ProvideStripeService(container *dig.Container) error {
	if e := container.Provide(
		func(
			config *config.Config,
			encryptService *services.EncryptService) *stripe.ClientStripeProvider {
			return stripe.NewClientStripeProvider(config, encryptService)
		}); e != nil {
		return e
	}
	if e := container.Provide(
		func(clientProvider *stripe.ClientStripeProvider) stripeInterface.AccountStripeServiceInterface {
			stripeImpl := stripe.NewAccountStripeService(clientProvider.MainClient())
			return stripeImpl
		}); e != nil {
		return e
	}
	if e := container.Provide(
		func(clientProvider *stripe.ClientStripeProvider) stripeInterface.CardStripeServiceInterface {
			return stripe.NewCardStripeService(clientProvider.MainClient())
		}); e != nil {
		return e
	}
	return nil
}
