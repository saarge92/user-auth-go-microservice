package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blockloop/scan"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	customErrors "go-user-microservice/internal/pkg/errors"
)

type CurrencyRepository struct {
	databaseInstance database.Database
}

func NewCurrencyRepository(databaseInstance database.Database) *CurrencyRepository {
	return &CurrencyRepository{databaseInstance: databaseInstance}
}

func (r *CurrencyRepository) GetByCode(ctx context.Context, code string) (*entites.Currency, error) {
	query := `SELECT * FROM currencies WHERE code = ?`
	currency := &entites.Currency{}

	currencyRow, currencyError := r.databaseInstance.QueryContext(ctx, query, code)
	if currencyError != nil {
		return nil, customErrors.DatabaseError(currencyError)
	}

	if e := scan.Row(currency, currencyRow); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, errorlists.CurrencyNotFoundErr
		}
		return nil, customErrors.DatabaseError(e)
	}

	return currency, nil
}
