package repositories

import (
	"go-user-microservice/internal/pkg/entites"
)

type CurrencyRepositoryInterface interface {
	GetByCode(code string) (*entites.Currency, error)
}
