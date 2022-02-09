package services

import (
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type StripeChargeServiceMock struct{}

func (s *StripeChargeServiceMock) CardCharge(cardInfo *dto.StripeCardCustomerChargeCreate) (*stripe.Charge, error) {
	return &stripe.Charge{
		ID: uuid.New().String(),
	}, nil
}
