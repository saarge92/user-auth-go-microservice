package repositories

import (
	"context"
	"go-user-microservice/internal/pkg/entites"
)

type CountryRepository interface {
	GetByCodeTwo(ctx context.Context, code string) (*entites.Country, error)
}
