package card

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/repositories"
)

type RepositoryCard struct {
	db *sqlx.DB
}

func NewRepositoryCard(db *sqlx.DB) *RepositoryCard {
	return &RepositoryCard{db: db}
}

func (r *RepositoryCard) Create(ctx context.Context, card *Card) error {
	query := `INSERT INTO cards (user_id, is_default, number, expire_day, expire_month, expire_year)
				VALUES (:user_id, :is_default, :number, :expire_day, :expire_month, :expire_year)`
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
