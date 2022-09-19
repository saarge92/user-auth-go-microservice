package card

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcServerCard struct {
	cardService        *ServiceCard
	transactionHandler *database.TransactionHandlerDB
}

func NewGrpcServerCard(cardService *ServiceCard, transactionHandler *database.TransactionHandlerDB) *GrpcServerCard {
	return &GrpcServerCard{
		cardService:        cardService,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcServerCard) CreateCard(ctx context.Context, request *core.CreateCardRequest) (*core.CreateCardResponse, error) {
	cardForm := forms.CreateCard{CreateCardRequest: request}
	if e := cardForm.Validate(); e != nil {
		return nil, e
	}

	transactionHandler := database.NewTypedTransaction[*core.CreateCardResponse](s.transactionHandler)

	return transactionHandler.WithCtx(ctx, func(ctx context.Context) (*core.CreateCardResponse, error) {
		cardInfo, e := s.cardService.Create(ctx, cardForm)
		if e != nil {
			return nil, e
		}
		return &core.CreateCardResponse{ExternalId: cardInfo.ExternalID}, nil
	})
}

func (s *GrpcServerCard) MyCards(ctx context.Context, _ *empty.Empty) (*core.MyCardsResponse, error) {
	cards, e := s.cardService.MyCards(ctx)
	if e != nil {
		return nil, e
	}
	return CardsListToGrpc(cards), nil
}
