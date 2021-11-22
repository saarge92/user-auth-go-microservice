package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/card/domain"
	entities2 "go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	stripeServices "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceCard struct {
	cardRepository    domain.CardRepositoryInterface
	cardStripeService stripeServices.CardStripeServiceInterface
}

func NewServiceCard(
	cardRepository domain.CardRepositoryInterface,
	cardStripeService stripeServices.CardStripeServiceInterface,
) *ServiceCard {
	return &ServiceCard{
		cardRepository:    cardRepository,
		cardStripeService: cardStripeService,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm *forms.CreateCard) (*entities2.Card, error) {
	var user *entities.User
	var convertOk bool
	if user, convertOk = ctx.Value(dictionary.User).(*entities.User); !convertOk {
		return nil, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}
	cardStripeDto := &dto.StripeCardCreate{
		Number:                 cardForm.CardNumber,
		ExpireMonth:            uint8(cardForm.ExpireMonth),
		ExpireYear:             cardForm.ExpireYear,
		CVC:                    cardForm.Cvc,
		StripePaymentAccountID: user.ProviderPaymentID,
	}
	var cardStripe *stripe.Card
	var cardError error
	cardChannelSync := make(chan interface{})
	go func() {
		cardStripe, cardError = s.cardStripeService.CreateCard(cardStripeDto, cardChannelSync)
	}()
	<-cardChannelSync
	if cardError != nil {
		return nil, cardError
	}
	cardEntity := &entities2.Card{}
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = user.ID
	cardEntity.ExternalID = uuid.New().String()
	cardEntity.ExternalProviderID = cardStripe.ID
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}
	return cardEntity, nil
}
