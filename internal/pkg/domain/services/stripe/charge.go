package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type ChargeService interface {
	CardCharge(cardInfo *dto.StripeCardCustomerChargeCreate) (*stripe.Charge, error)
}
