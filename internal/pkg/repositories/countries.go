package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blockloop/scan"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	sharedErrors "go-user-microservice/internal/pkg/errors"
)

type CountryRepository struct {
	databaseInstance database.Database
}

func NewCountryRepository(databaseInstance database.Database) *CountryRepository {
	return &CountryRepository{
		databaseInstance: databaseInstance,
	}
}

func (r *CountryRepository) GetByCodeTwo(ctx context.Context, code string) (*entites.Country, error) {
	query := `SELECT * from countries WHERE code_2 = ?`
	country := &entites.Country{}

	countryRow, countryError := r.databaseInstance.QueryContext(ctx, query, code)
	if countryError != nil {
		return nil, sharedErrors.DatabaseError(countryError)
	}

	if e := scan.Row(country, countryRow); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, errorlists.CountryNotFoundErr
		}
		return nil, sharedErrors.DatabaseError(e)
	}

	return country, nil
}
