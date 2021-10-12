package repositories

import "go-user-microservice/internal/entites"

type CurrencyRepositoryInterface interface {
	GetByCode(code string) (*entites.Currency, error)
}
