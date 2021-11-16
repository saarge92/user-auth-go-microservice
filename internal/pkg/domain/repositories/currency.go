package repositories

import (
	"context"
	"go-user-microservice/internal/pkg/entites"
)

type CurrencyRepositoryInterface interface {
	GetByCode(ctx context.Context, code string) (*entites.Currency, error)
}
