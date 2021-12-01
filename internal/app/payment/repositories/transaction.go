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

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, record *entities.Transaction) error {
	record.ExternalID = uuid.New().String()
	record.CreatedAt = time.Now()
	query := `INSERT INTO transactions
				(external_id, transaction_type, external_provider_id, 
				from_user_id, to_user_id, amount, created_at) 
			VALUES (:external_id, :transaction_type, :external_provider_id, 
					:from_user_id, :to_user_id, :amount, :created_at)`
	var dbError error
	var result sql.Result
	tx := repositories.GetDBTransaction(ctx)
	if tx != nil {
		result, dbError = tx.NamedExec(query, record)
	} else {
		result, dbError = r.db.NamedExec(query, record)
	}
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	record.ID = uint64(repositories.LastInsertID(result))
	return nil
}
