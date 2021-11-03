package stripe

import (
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/services"
)

type ClientStripeWrapper struct {
	client         *client.API
	encryptService *services.EncryptService
}

func NewClientStripe(
	config *config.Config,
	encryptService *services.EncryptService,
) *ClientStripeWrapper {
	newClientStripeWrapper := &ClientStripeWrapper{
		encryptService: encryptService,
		client:         &client.API{},
	}
	secretKey, e := encryptService.Decrypt([]byte(config.SecretStripeKey), []byte(config.SecretEncryptionKey))
	if e != nil {
		log.Error(e)
	}
	newClientStripeWrapper.client.Init(string(secretKey), nil)
	return newClientStripeWrapper
}
