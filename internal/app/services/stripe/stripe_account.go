package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/dto"
)

type AccountStripeService struct {
	clientStripe *ClientStripeWrapper
}

func NewAccountStripeService(client *ClientStripeWrapper) *AccountStripeService {
	return &AccountStripeService{
		clientStripe: client,
	}
}

func (s *AccountStripeService) Create(data *dto.StripeAccountCreate) (*stripe.Account, error) {
	accountType := "custom"
	accountParams := &stripe.AccountParams{
		Email: &data.Email,
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(data.CardPaymentsRequested),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(data.TransferRequested),
			},
		},
		Type: &accountType,
	}
	if data.Country != "" {
		accountParams.Country = &data.Country
	}
	account, e := s.clientStripe.client.Account.New(accountParams)
	if e != nil {
		return nil, e
	}
	return account, e
}
