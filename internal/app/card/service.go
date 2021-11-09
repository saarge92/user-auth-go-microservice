package card

import (
	"context"
	"fmt"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/services/stripe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceCard struct {
	cardRepository    *RepositoryCard
	cardStripeService *stripe.CardStripeService
}

func NewServiceCard(
	cardRepository *RepositoryCard,
	cardStripeService *stripe.CardStripeService,
) *ServiceCard {
	return &ServiceCard{
		cardRepository:    cardRepository,
		cardStripeService: cardStripeService,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm *forms.CreateCard) (*Card, error) {
	var user *entities.User
	var convertOk bool
	if user, convertOk = ctx.Value(dictionary.User).(*entities.User); !convertOk {
		return nil, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}
	cardStripeDto := &dto.StripeCardCreate{
		Number:      cardForm.CardNumber,
		ExpireMonth: uint8(cardForm.ExpireMonth),
		ExpireYear:  cardForm.ExpireYear,
		CVC:         cardForm.Cvc,
	}
	cardStripe, e := s.cardStripeService.CreateCard(cardStripeDto)
	if e != nil {
		return nil, e
	}
	cardEntity := &Card{}
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = user.ID
	cardEntity.ExternalProviderID = cardStripe.ID
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}
	return cardEntity, nil
}
