package services

import (
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type AccountStripeServiceMock struct{}

func (s *AccountStripeServiceMock) Create(data *dto.StripeAccountCreate) (*stripe.Account, *stripe.Customer, error) {
	account := &stripe.Account{
		ID:           uuid.New().String(),
		Email:        data.Email,
		Country:      data.Country,
		Capabilities: nil,
	}
	customer := &stripe.Customer{
		ID: uuid.New().String(),
	}
	return account, customer, nil
}
