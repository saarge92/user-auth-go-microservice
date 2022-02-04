package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/dto"
)

type AccountStripeService struct {
	mainStripeClient *client.API
}

func NewAccountStripeService(client *client.API) *AccountStripeService {
	return &AccountStripeService{
		mainStripeClient: client,
	}
}

func (s *AccountStripeService) Create(data *dto.StripeAccountCreate) (*stripe.Account, *stripe.Customer, error) {
	accountType := dictionary.CustomStripeAccountType
	accountParams := &stripe.AccountParams{
		Email: &data.Email,
		Type:  (*string)(&accountType),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(data.CardPayments),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(data.Transfers),
			},
		},
	}
	if data.Country != "" {
		accountParams.Country = &data.Country
	}
	account, e := s.mainStripeClient.Account.New(accountParams)
	if e != nil {
		return nil, nil, e
	}
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(data.Email),
	}
	customer, e := s.mainStripeClient.Customers.New(customerParams)
	if e != nil {
		return nil, nil, e
	}
	return account, customer, nil
}
