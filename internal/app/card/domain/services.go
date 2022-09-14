package domain

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type StripeCardService interface {
	CreateCard(cardData dto.StripeCardCreate) (*stripe.Card, error)
}
