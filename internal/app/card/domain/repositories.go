package domain

import (
	"context"
	"go-user-microservice/internal/app/card/entities"
)

type CardRepositoryInterface interface {
	Create(ctx context.Context, card *entities.Card) error
}
