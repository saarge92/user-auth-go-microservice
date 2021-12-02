package card

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/repositories"
	"time"
)

type RepositoryCard struct {
	db *sqlx.DB
}

func NewRepositoryCard(db *sqlx.DB) *RepositoryCard {
	return &RepositoryCard{db: db}
}

func (r *RepositoryCard) Create(ctx context.Context, card *entities.Card) error {
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now
	query := `INSERT INTO cards (
                user_id, is_default, number, external_provider_id, external_id,
                expire_month, expire_year, created_at, updated_at)
				VALUES (:user_id, :is_default, :number, :external_provider_id, :external_id,
				     	:expire_month, :expire_year, :created_at, :updated_at)`
	tx := repositories.GetDBTransaction(ctx)
	var result sql.Result
	var dbError error
	if tx != nil {
		result, dbError = tx.NamedExec(query, card)
	} else {
		result, dbError = r.db.NamedExec(query, card)
	}
	if dbError != nil {
		return errors.DatabaseError(dbError)
	}
	card.ID = uint64(repositories.LastInsertID(result))
	return nil
}

func (r *RepositoryCard) ListByCardID(ctx context.Context, userID uint64) ([]entities.Card, error) {
	query := `SELECT * FROM cards WHERE user_id = ?`
	var cards []entities.Card
	tx := repositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Select(&cards, query, userID)
	} else {
		dbError = r.db.Select(&cards, query, userID)
	}
	if dbError != nil {
		return nil, errors.DatabaseError(dbError)
	}
	return cards, dbError
}

func (r *RepositoryCard) OneByCardAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*entities.Card, error) {
	query := `SELECT * FROM cards WHERE external_id = ? AND user_id = ?`
	card := new(entities.Card)
	tx := repositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Get(card, query, externalID, userID)
	} else {
		dbError = r.db.Get(card, query, externalID, userID)
	}
	if dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.DatabaseError(dbError)
	}
	return card, nil
}
