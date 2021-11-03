package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/dto"
)

type AccountStripeServiceInterface interface {
	Create(data *dto.StripeAccountCreate) (*stripe.Account, error)
}
