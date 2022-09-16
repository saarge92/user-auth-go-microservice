package domain

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

//go:generate go run github.com/vektra/mockery/v2@latest --with-expecter --case underscore --all --output=./../mocks

type StripeCardService interface {
	CreateCard(cardData dto.StripeCardCreate) (*stripe.Card, error)
}

type StripeBackend interface {
	stripe.Backend
}
