package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/payment/entities"
	customErrors "go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/repositories"
	"time"
)

type OperationStoryRepository struct {
	db *sqlx.DB
}

func NewOperationStoryRepository(db *sqlx.DB) *OperationStoryRepository {
	return &OperationStoryRepository{
		db: db,
	}
}

func (r *OperationStoryRepository) Create(ctx context.Context, operationStory *entities.OperationStory) error {
	operationStory.ExternalID = uuid.New().String()
	operationStory.CreatedAt = time.Now()
	query := ` INSERT INTO operation_stories 
     			(
     			 	external_id, user_id, card_id, amount, balance_before, balance_after,
					external_provider_id, operation_type_id, created_at
				)
				VALUES 
				(
					:external_id, :user_id, :card_id, :amount, :balance_before, :balance_after,
					:external_provider_id, :operation_type_id, :created_at
				)
				`
	tx := repositories.GetDBTransaction(ctx)
	var dbError error
	var result sql.Result
	if tx != nil {
		result, dbError = tx.NamedExec(query, operationStory)
	} else {
		result, dbError = r.db.NamedExec(query, operationStory)
	}
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	operationStory.ID = uint64(repositories.LastInsertID(result))
	return nil
}
