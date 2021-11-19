package providers

import (
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/services"
)

type ClientStripeProvider struct {
	mainClient *client.API
}

func NewClientStripeProvider(
	config *config.Config,
) *ClientStripeProvider {
	encryptService := &services.EncryptService{}
	newClientStripeProvider := &ClientStripeProvider{
		mainClient: &client.API{},
	}
	secretKey, e := encryptService.Decrypt([]byte(config.SecretStripeKey), []byte(config.SecretEncryptionKey))
	if e != nil {
		log.Error(e)
	}
	newClientStripeProvider.mainClient.Init(string(secretKey), nil)
	return newClientStripeProvider
}

func (p *ClientStripeProvider) MainClient() *client.API {
	return p.mainClient
}
