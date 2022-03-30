package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/entites"
	customErrors "go-user-microservice/internal/pkg/errors"
)

type CurrencyRepository struct {
	db *sqlx.DB
}

func NewCurrencyRepository(db *sqlx.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*entites.Currency, error) {
	dbConn := db.GetDBConnection(ctx, r.db)

	query := `SELECT * FROM currencies WHERE code = ?`
	currency := &entites.Currency{}

	if dbError := dbConn.Get(currency, query, code); dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, nil
		}
		return nil, customErrors.DatabaseError(dbError)
	}
	return currency, nil
}
