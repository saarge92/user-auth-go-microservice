package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/errors"
	"google.golang.org/grpc/codes"
)

type CountryRepository struct {
	db *sqlx.DB
}

func NewCountryRepository(db *sqlx.DB) *CountryRepository {
	return &CountryRepository{
		db: db,
	}
}

func (r *CountryRepository) GetByCodeTwo(ctx context.Context, code string) (*entites.Country, error) {
	query := `SELECT * from countries WHERE code_2 = ?`
	tx := getDBConnection(ctx, r.db)
	country := &entites.Country{}
	e := tx.Get(country, query, code)
	if e != nil {
		if e == sql.ErrNoRows {
			return nil, errors.CustomDatabaseError(codes.NotFound, errorlists.CountryNotFound)
		}
		return nil, e
	}
	e = tx.Commit()
	if e != nil {
		return nil, e
	}
	return country, nil
}
