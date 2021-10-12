package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errors"
)

type CurrencyRepository struct {
	db *sqlx.DB
}

func NewCurrencyRepository(db *sqlx.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) GetByCode(code string) (*entites.Currency, error) {
	query := `SELECT * FROM currencies WHERE code = ?`
	currency := &entites.Currency{}
	e := r.db.Select(currency, query, code)
	if e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.DatabaseError(e)
	}
	return currency, nil
}
