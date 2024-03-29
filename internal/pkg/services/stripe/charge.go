package stripe

import (
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/dto"
)

type ChargeStripeService struct {
	stripeClient *client.API
}

func NewChargeStripeService(
	stripeClient *client.API,
) *ChargeStripeService {
	return &ChargeStripeService{
		stripeClient: stripeClient,
	}
}

func (s *ChargeStripeService) CardCharge(cardInfo dto.StripeCardCustomerChargeCreate) (*stripe.Charge, error) {
	amount := cardInfo.Amount.Mul(decimal.NewFromInt(100)).IntPart()
	cardChargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(cardInfo.Currency),
		Description: stripe.String(cardInfo.Description),
		Customer:    stripe.String(cardInfo.CustomerID),
		Source: &stripe.SourceParams{
			Card: &stripe.CardParams{
				ID: cardInfo.CardID,
			},
		},
	}
	chargeResponse, e := s.stripeClient.Charges.New(cardChargeParams)
	if e != nil {
		return nil, e
	}
	return chargeResponse, e
}
