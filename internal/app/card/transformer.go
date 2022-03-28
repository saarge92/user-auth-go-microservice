package card

import (
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/pkg/protobuf/core"
)

func CardsListToGrpc(cards []entities.Card) *core.MyCardsResponse {
	grpcCards := make([]*core.Card, 0, len(cards))
	for _, cardElement := range cards {
		grpcCardElement := &core.Card{
			CardNumber:  cardElement.Number,
			ExternalId:  cardElement.ExternalID,
			ExpireYear:  cardElement.ExpireYear,
			ExpireMonth: cardElement.ExpireMonth,
		}
		grpcCards = append(grpcCards, grpcCardElement)
	}
	return &core.MyCardsResponse{
		Cards: grpcCards,
	}
}
