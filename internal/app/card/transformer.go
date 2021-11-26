package card

import (
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/pkg/protobuf/card"
)

func CardsListToGrpc(cards []entities.Card) *card.MyCardsResponse {
	grpcCards := make([]*card.Card, 0, len(cards))
	for _, cardElement := range cards {
		grpcCardElement := &card.Card{
			CardNumber:  cardElement.Number,
			ExternalId:  cardElement.ExternalID,
			ExpireYear:  cardElement.ExpireYear,
			ExpireMonth: cardElement.ExpireMonth,
		}
		grpcCards = append(grpcCards, grpcCardElement)
	}
	return &card.MyCardsResponse{
		Cards: grpcCards,
	}
}
