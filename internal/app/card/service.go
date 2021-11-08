package card

import (
	"context"
	"fmt"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceCard struct {
	cardRepository *RepositoryCard
}

func NewServiceCard(
	cardRepository *RepositoryCard,
) *ServiceCard {
	return &ServiceCard{
		cardRepository: cardRepository,
	}
}

func (s *ServiceCard) Create(ctx context.Context, cardForm *forms.CreateCard) (*Card, error) {
	var user *entities.User
	var convertOk bool
	if user, convertOk = ctx.Value(dictionary.User).(*entities.User); !convertOk {
		return nil, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}
	cardEntity := &Card{}
	cardEntity.Number = cardForm.CardNumber
	cardEntity.ExpireDay = cardForm.ExpireDay
	cardEntity.ExpireMonth = cardForm.ExpireMonth
	cardEntity.ExpireYear = cardForm.ExpireYear
	cardEntity.IsDefault = cardForm.IsDefault
	cardEntity.UserID = user.ID
	if e := s.cardRepository.Create(ctx, cardEntity); e != nil {
		return nil, e
	}
	return cardEntity, nil
}
