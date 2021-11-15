package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/dto"
)

type AccountStripeService struct {
	clientStripe *client.API
}

func NewAccountStripeService(client *client.API) *AccountStripeService {
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
	account, e := s.clientStripe.Account.New(accountParams)
	if e != nil {
		return nil, e
	}
	return account, e
}
