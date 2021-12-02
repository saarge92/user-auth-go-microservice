package repositories

import (
	"context"
	"go-user-microservice/internal/pkg/entites"
)

type CurrencyRepository interface {
	GetByCode(ctx context.Context, code string) (*entites.Currency, error)
}
