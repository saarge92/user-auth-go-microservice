package card

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type RepositoryCard struct {
	db *sqlx.DB
}

func NewRepositoryCard(db *sqlx.DB) *RepositoryCard {
	return &RepositoryCard{db: db}
}

func (r *RepositoryCard) Create(ctx context.Context, card *entities.Card) error {
	dbConn := db.GetDBConnection(ctx, r.db)
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now
	query := `INSERT INTO cards (
                user_id, is_default, number, external_provider_id, external_id,
                expire_month, expire_year, created_at, updated_at)
				VALUES (:user_id, :is_default, :number, :external_provider_id, :external_id,
				     	:expire_month, :expire_year, :created_at, :updated_at)`

	var result sql.Result
	var dbError error
	result, dbError = dbConn.NamedExec(query, card)
	if dbError != nil {
		return dbError
	}
	card.ID = uint64(db.LastInsertID(result))
	return nil
}

func (r *RepositoryCard) ListByCardID(ctx context.Context, userID uint64) ([]entities.Card, error) {
	dbConn := db.GetDBConnection(ctx, r.db)

	query := `SELECT * FROM cards WHERE user_id = ?`
	var cards []entities.Card
	if e := dbConn.Select(&cards, query, userID); e != nil {
		return nil, e
	}

	return cards, nil
}

func (r *RepositoryCard) OneByCardAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*entities.Card, error) {
	dbConn := db.GetDBConnection(ctx, r.db)

	query := `SELECT * FROM cards WHERE external_id = ? AND user_id = ?`
	card := new(entities.Card)

	if e := dbConn.Get(card, query, externalID, userID); e != nil {
		if e == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, errorlists.CardNotFound)
		}
		return nil, errors.DatabaseError(e)
	}

	return card, nil
}
