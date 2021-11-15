package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"go-user-microservice/internal/pkg/dto"
	"strconv"
)

type CardStripeService struct {
	stripeClient *client.API
}

func NewCardStripeService(stripeClient *client.API) *CardStripeService {
	return &CardStripeService{
		stripeClient: stripeClient,
	}
}

func (s *CardStripeService) CreateCard(
	cardData *dto.StripeCardCreate,
	syncChannel chan interface{},
) (*stripe.Card, error) {
	defer close(syncChannel)
	expireMonth := strconv.Itoa(int(cardData.ExpireMonth))
	expireYear := strconv.Itoa(int(cardData.ExpireYear))
	cvc := strconv.Itoa(int(cardData.CVC))
	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(cardData.Number),
			ExpMonth: stripe.String(expireMonth),
			ExpYear:  stripe.String(expireYear),
			CVC:      stripe.String(cvc),
			Currency: stripe.String("USD"),
		},
	}
	token, e := s.stripeClient.Tokens.New(tokenParams)
	if e != nil {
		return nil, e
	}
	cardParams := &stripe.CardParams{
		Token:   stripe.String(token.ID),
		Account: stripe.String(cardData.StripePaymentAccountID),
	}
	cardStripe, e := s.stripeClient.Cards.New(cardParams)
	if e != nil {
		return nil, e
	}
	return cardStripe, nil
}
