package domain

import (
	"context"
	"go-user-microservice/internal/app/card"
)

type CardRepositoryInterface interface {
	Create(ctx context.Context, card *card.Card) error
}
