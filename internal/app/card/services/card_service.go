package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/card/domain"
	cardEntities "go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/card/errors"
	"go-user-microservice/internal/app/card/forms"
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/pkg/dictionary"
	stripeServices "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceCard struct {
	cardRepository    domain.CardRepository
	cardStripeService stripeServices.CardStripeService
}

func NewServiceCard(
	cardRepository domain.CardRepository,
	cardStripeService stripeServices.CardStripeService,
) *ServiceCard {
	return &ServiceCard{
		cardRepository:    cardRepository,
		cardStripeService: cardStripeService,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm *forms.CreateCard) (*cardEntities.Card, error) {
	cardExist, cardErr := s.cardRepository.ExistByCardNumber(ctx, cardForm.CardNumber)
	if cardErr != nil {
		return nil, cardErr
	}
	if cardExist {
		return nil, errors.ErrCardNotFound
	}

	var userRole *userDto.UserRole
	var convertOk bool
	if userRole, convertOk = ctx.Value(dictionary.User).(*userDto.UserRole); !convertOk {
		return nil, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}
	cardStripeDto := &dto.StripeCardCreate{
		Number:             cardForm.CardNumber,
		ExpireMonth:        uint8(cardForm.ExpireMonth),
		ExpireYear:         cardForm.ExpireYear,
		CVC:                cardForm.Cvc,
		AccountProviderID:  userRole.User.AccountProviderID,
		CustomerProviderID: userRole.User.CustomerProviderID,
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
	cardEntity := &cardEntities.Card{}
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = userRole.User.ID
	cardEntity.ExternalID = uuid.New().String()
	cardEntity.ExternalProviderID = cardStripe.ID
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}
	return cardEntity, nil
}

func (s *ServiceCard) MyCards(
	ctx context.Context,
) ([]cardEntities.Card, error) {
	var user *userDto.UserRole
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*userDto.UserRole); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	cards, e := s.cardRepository.ListByCardID(ctx, user.User.ID)
	if e != nil {
		return nil, e
	}
	return cards, nil
}
