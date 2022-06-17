package card

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Knetic/go-namedParameterQuery"
	"github.com/blockloop/scan"
	"go-user-microservice/internal/app/card/entities"
	errors2 "go-user-microservice/internal/app/card/errors"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/db"
	"time"
)

type RepositoryCard struct {
	databaseInstance database.Database
}

func NewRepositoryCard(databaseInstance database.Database) *RepositoryCard {
	return &RepositoryCard{databaseInstance: databaseInstance}
}

func (r *RepositoryCard) Create(ctx context.Context, card *entities.Card) error {
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now
	query := `INSERT INTO cards (
                user_id, is_default, number, external_provider_id, external_id,
                expire_month, expire_year, created_at, updated_at)
				VALUES (:userId, :isDefault, :number, :externalProviderId, :externalId,
				     	:expireMonth, :expireYear, :createdAt, :updatedAt)`

	queryNamed := namedParameterQuery.NewNamedParameterQuery(query)
	insertParams := map[string]interface{}{
		"userId":             card.UserID,
		"isDefault":          card.IsDefault,
		"number":             card.Number,
		"externalProviderId": card.ExternalProviderID,
		"externalId":         card.ExternalID,
		"expireMonth":        card.ExpireMonth,
		"expireYear":         card.ExpireYear,
		"createdAt":          card.CreatedAt,
		"updatedAt":          card.UpdatedAt,
	}
	queryNamed.SetValuesFromMap(insertParams)
	result, dbError := r.databaseInstance.ExecContext(ctx, queryNamed.GetParsedQuery(), queryNamed.GetParsedParameters()...)
	if dbError != nil {
		return dbError
	}

	card.ID = uint64(db.LastInsertID(result))
	return nil
}

func (r *RepositoryCard) ListByCardID(ctx context.Context, userID uint64) ([]entities.Card, error) {
	query := `SELECT * FROM cards WHERE user_id = ?`
	var cards []entities.Card

	cardRows, cardError := r.databaseInstance.QueryContext(ctx, query, userID)
	if cardError != nil {
		return nil, cardError
	}

	if e := scan.Rows(&cards, cardRows); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, e
	}

	return cards, nil
}

func (r *RepositoryCard) OneByCardAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*entities.Card, error) {
	query := `SELECT * FROM cards WHERE external_id = ? AND user_id = ?`
	card := new(entities.Card)

	cardRow, cardError := r.databaseInstance.QueryContext(ctx, query, externalID, userID)
	if cardError != nil {
		return nil, cardError
	}

	if e := scan.Row(card, cardRow); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, errors2.ErrCardNotFound
		}
	}

	return card, nil
}

func (r *RepositoryCard) ExistByCardNumber(ctx context.Context, cardNumber string) (bool, error) {
	query := "SELECT COUNT(*) > 0 FROM cards WHERE number = ?"
	var exist bool
	if e := r.databaseInstance.QueryRowContext(ctx, query, cardNumber).Scan(&exist); e != nil {
		if !errors.Is(e, sql.ErrNoRows) {
			return false, e
		}
		return false, nil
	}

	return exist, nil
}
