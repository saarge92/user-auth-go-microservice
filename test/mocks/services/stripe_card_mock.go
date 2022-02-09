package services

import (
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type StripeCardServiceMock struct{}

func (s StripeCardServiceMock) CreateCard(cardData *dto.StripeCardCreate, syncChannel chan interface{}) (*stripe.Card, error) {
	defer close(syncChannel)
	return &stripe.Card{
		ID: uuid.New().String(),
	}, nil
}
