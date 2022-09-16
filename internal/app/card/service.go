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
)

type ServiceCard struct {
	cardRepository    domain.CardRepository
	cardStripeService domain.StripeCardService
}

func NewServiceCard(cardRepository domain.CardRepository, cardStripeService domain.StripeCardService) *ServiceCard {
	return &ServiceCard{
		cardRepository:    cardRepository,
		cardStripeService: cardStripeService,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm forms.CreateCard) (*entities.Card, error) {
	if cardErr := s.checkCardNumberExist(ctx, cardForm.CardNumber); cardErr != nil {
		return nil, cardErr
	}

	userRoleDto, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}

	cardStripeDto := dto.StripeCardCreate{
		Number:            cardForm.CardNumber,
		ExpireMonth:       uint8(cardForm.ExpireMonth),
		ExpireYear:        cardForm.ExpireYear,
		CVC:               cardForm.Cvc,
		AccountProviderID: userRoleDto.User.AccountProviderID,
	}
	cardStripe, cardError := s.cardStripeService.CreateCard(cardStripeDto)
	if cardError != nil {
		return nil, cardError
	}

	return s.initCardRecord(ctx, cardForm, userRoleDto.User.ID, cardStripe.ID)
}

func (s *ServiceCard) MyCards(ctx context.Context) ([]entities.Card, error) {
	userRoleDto, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}
	cards, e := s.cardRepository.ListByCardID(ctx, userRoleDto.User.ID)
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

func (s *ServiceCard) initCardRecord(ctx context.Context, cardForm forms.CreateCard, userID uint64, providerID string) (*entities.Card, error) {
	cardEntity := &entities.Card{}
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
