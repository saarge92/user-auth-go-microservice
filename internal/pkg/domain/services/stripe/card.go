package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
)

type CardStripeServiceInterface interface {
	CreateCard(cardData *dto.StripeCardCreate, syncChannel chan interface{}) (*stripe.Card, error)
}
