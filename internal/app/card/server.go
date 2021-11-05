package card

import (
	"context"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/pkg/protobuf/card"
)

type GrpcServerCard struct {
	cardService     *ServiceCard
	cardFormBuilder *forms.CardFormBuilder
}

func NewGrpcServerCard(
	cardFormBuilder *forms.CardFormBuilder,
	cardService *ServiceCard,
) *GrpcServerCard {
	return &GrpcServerCard{
		cardService:     cardService,
		cardFormBuilder: cardFormBuilder,
	}
}

func (s *GrpcServerCard) CreateCard(
	ctx context.Context,
	request *card.CreateCardRequest,
) (*card.CreateCardResponse, error) {
	cardForm := s.cardFormBuilder.CreateCreateForm(request)
	if e := cardForm.Validate(); e != nil {
		return nil, e
	}
	_, e := s.cardService.Create(ctx, cardForm)
	if e != nil {
		return nil, e
	}
	return &card.CreateCardResponse{ExternalId: ""}, nil
}
