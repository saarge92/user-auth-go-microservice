package card

import (
	"context"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/card/services"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcServerCard struct {
	cardService        *services.ServiceCard
	cardFormBuilder    *forms.CardFormBuilder
	transactionHandler *db.TransactionHandlerDB
}

func NewGrpcServerCard(
	cardFormBuilder *forms.CardFormBuilder,
	cardService *services.ServiceCard,
	transactionHandler *db.TransactionHandlerDB,
) *GrpcServerCard {
	return &GrpcServerCard{
		cardService:        cardService,
		cardFormBuilder:    cardFormBuilder,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcServerCard) CreateCard(
	ctx context.Context,
	request *core.CreateCardRequest,
) (response *core.CreateCardResponse, e error) {
	cardForm := s.cardFormBuilder.CreateCreateForm(request)
	if e = cardForm.Validate(); e != nil {
		return
	}

	ctx, tx, e := s.transactionHandler.Create(ctx, nil)
	if e != nil {
		return
	}

	defer func() {
		e = db.HandleTransaction(tx, e)
	}()

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
