package providers

import (
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/config"
)

type ClientStripeProvider struct {
	mainClient *client.API
}

func NewClientStripeProvider(
	config *config.Config,
) *ClientStripeProvider {
	newClientStripeProvider := &ClientStripeProvider{
		mainClient: &client.API{},
	}
	newClientStripeProvider.mainClient.Init(config.SecretStripeKey, nil)
	return newClientStripeProvider
}

func (p *ClientStripeProvider) MainClient() *client.API {
	return p.mainClient
}
