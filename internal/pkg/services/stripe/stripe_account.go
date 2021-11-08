package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/pkg/dictionary"
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
	accountType := dictionary.StandardStripeAccountType
	accountParams := &stripe.AccountParams{
		Email: &data.Email,
		Type:  (*string)(&accountType),
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
