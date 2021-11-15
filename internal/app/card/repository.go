package card

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
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

func (r *RepositoryCard) Create(ctx context.Context, card *Card) error {
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
