package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/pkg/dto"
	"strconv"
)

type CardStripeService struct {
	client *ClientStripeWrapper
}

func NewCardStripeService(client *ClientStripeWrapper) *CardStripeService {
	return &CardStripeService{
		client: client,
	}
}

func (s *CardStripeService) CreateCard(cardData *dto.StripeCardCreate) (*stripe.Card, error) {
	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(cardData.Number),
			ExpMonth: stripe.String(string(cardData.ExpireMonth)),
			ExpYear:  stripe.String(string(cardData.ExpireMonth)),
			CVC:      stripe.String(strconv.Itoa(int(cardData.CVC))),
		},
	}
	token, e := s.client.client.Tokens.New(tokenParams)
	if e != nil {
		return nil, e
	}
	cardParams := &stripe.CardParams{
		Token:   stripe.String(token.ID),
		Account: stripe.String(cardData.StripePaymentAccountID),
	}
	cardStripe, e := s.client.client.Cards.New(cardParams)
	if e != nil {
		return nil, e
	}
	return cardStripe, nil
}
