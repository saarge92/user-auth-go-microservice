package services

import (
	"context"
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
	"go-user-microservice/internal/pkg/grpc"
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
	if cardErr := s.checkCardNumberExist(ctx, cardForm.CardNumber); cardErr != nil {
		return nil, cardErr
	}

	userRoleDto, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}

	cardStripeDto := &dto.StripeCardCreate{
		Number:             cardForm.CardNumber,
		ExpireMonth:        uint8(cardForm.ExpireMonth),
		ExpireYear:         cardForm.ExpireYear,
		CVC:                cardForm.Cvc,
		AccountProviderID:  userRoleDto.User.AccountProviderID,
		CustomerProviderID: userRoleDto.User.CustomerProviderID,
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

	return s.initCardRecord(ctx, cardForm, userRoleDto.User.ID, cardStripe.ID)
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

func (s *ServiceCard) checkCardNumberExist(ctx context.Context, cardNumber string) error {
	cardExist, cardErr := s.cardRepository.ExistByCardNumber(ctx, cardNumber)
	if cardErr != nil {
		return cardErr
	}
	if cardExist {
		return errors.ErrCardNotFound
	}

	return nil
}

func (s *ServiceCard) initCardRecord(ctx context.Context, cardForm *forms.CreateCard, userID uint64, providerID string) (*cardEntities.Card, error) {
	cardEntity := &cardEntities.Card{}
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = userID
	cardEntity.ExternalID = uuid.New().String()
	cardEntity.ExternalProviderID = providerID
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}

	return cardEntity, nil
}
