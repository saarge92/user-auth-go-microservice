package card

import (
	"context"
	"github.com/google/uuid"
	"go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/card/errors"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/grpc"
	"time"
)

type ServiceCard struct {
	cardRepository    domain.CardRepository
	cardStripeService domain.StripeCardService
}

type cardInitParam struct {
	userID     uint64
	providerID string
	cardForm   forms.CreateCard
}

func NewServiceCard(cardRepository domain.CardRepository, cardStripeService domain.StripeCardService) *ServiceCard {
	return &ServiceCard{
		cardRepository:    cardRepository,
		cardStripeService: cardStripeService,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm forms.CreateCard) (*entities.Card, error) {
	userData, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}

	if cardErr := s.checkCardNumberExist(ctx, cardForm.CardNumber); cardErr != nil {
		return nil, cardErr
	}

	cardStripeDto := dto.StripeCardCreate{
		Number:            cardForm.CardNumber,
		ExpireMonth:       uint8(cardForm.ExpireMonth),
		ExpireYear:        cardForm.ExpireYear,
		CVC:               cardForm.Cvc,
		AccountProviderID: userData.AccountProviderID,
	}
	cardStripe, cardError := s.cardStripeService.CreateCard(cardStripeDto)
	if cardError != nil {
		return nil, cardError
	}

	initCardParam := cardInitParam{
		cardForm:   cardForm,
		userID:     userData.ID,
		providerID: cardStripe.ID,
	}

	return s.initCardRecord(ctx, initCardParam)
}

func (s *ServiceCard) MyCards(ctx context.Context) ([]entities.Card, error) {
	userData, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}
	cards, e := s.cardRepository.ListByCardID(ctx, userData.ID)
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
		return errors.ErrCardAlreadyExist
	}

	return nil
}

func (s *ServiceCard) initCardRecord(ctx context.Context, initCardParam cardInitParam) (*entities.Card, error) {
	cardEntity := &entities.Card{}
	now := time.Now().Unix()
	cardForm := initCardParam.cardForm
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = initCardParam.userID
	cardEntity.ExternalID = uuid.New().String()
	cardEntity.ExternalProviderID = initCardParam.providerID
	cardEntity.CreatedAt = now
	cardEntity.UpdatedAt = now
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}

	return cardEntity, nil
}
