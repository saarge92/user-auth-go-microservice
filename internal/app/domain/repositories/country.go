package repositories

import (
	"context"
	"go-user-microservice/internal/app/entites"
)

type CountryRepositoryInterface interface {
	GetByCodeTwo(ctx context.Context, code string) (*entites.Country, error)
}
