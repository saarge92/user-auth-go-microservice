package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/errors"
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
	dbConn := db.GetDBConnection(ctx, r.db)

	query := `SELECT * from countries WHERE code_2 = ?`
	country := &entites.Country{}

	if dbError := dbConn.Get(country, query, code); dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, errors.CustomDatabaseError(codes.NotFound, errorlists.CountryNotFound)
		}
		return nil, errors.DatabaseError(dbError)
	}
	return country, nil
}
