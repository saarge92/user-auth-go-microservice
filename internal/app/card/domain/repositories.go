package domain

import (
	"context"
	"go-user-microservice/internal/app/card/entities"
)

type CardRepository interface {
	Create(ctx context.Context, cardEntity *entities.Card) error
	ListByCardID(ctx context.Context, userID uint64) ([]entities.Card, error)
	OneByCardAndUserID(ctx context.Context, externalID string, userID uint64) (*entities.Card, error)
	ExistByCardNumber(ctx context.Context, cardNumber string) (bool, error)
}
