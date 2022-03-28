package card

import (
	"context"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/card/services"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcServerCard struct {
	cardService     *services.ServiceCard
	cardFormBuilder *forms.CardFormBuilder
}

func NewGrpcServerCard(
	cardFormBuilder *forms.CardFormBuilder,
	cardService *services.ServiceCard,
) *GrpcServerCard {
	return &GrpcServerCard{
		cardService:     cardService,
		cardFormBuilder: cardFormBuilder,
	}
}

func (s *GrpcServerCard) CreateCard(
	ctx context.Context,
	request *core.CreateCardRequest,
) (*core.CreateCardResponse, error) {
	cardForm := s.cardFormBuilder.CreateCreateForm(request)
	if e := cardForm.Validate(); e != nil {
		return nil, e
	}
	cardInfo, e := s.cardService.Create(ctx, cardForm)
	if e != nil {
		return nil, e
	}
	return &core.CreateCardResponse{ExternalId: cardInfo.ExternalID}, nil
}

func (s *GrpcServerCard) MyCards(
	ctx context.Context,
	_ *core.MyCardsRequest,
) (*core.MyCardsResponse, error) {
	cards, e := s.cardService.MyCards(ctx)
	if e != nil {
		return nil, e
	}
	return CardsListToGrpc(cards), nil
}
