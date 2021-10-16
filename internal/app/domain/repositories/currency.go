package repositories

import (
	"go-user-microservice/internal/app/entites"
)

type CurrencyRepositoryInterface interface {
	GetByCode(code string) (*entites.Currency, error)
}
