package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type AccountStripeService interface {
	Create(data *dto.StripeAccountCreate) (*stripe.Account, *stripe.Customer, error)
}
