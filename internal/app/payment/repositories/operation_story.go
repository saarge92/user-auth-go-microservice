package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	"go-user-microservice/internal/pkg/db"
	customErrors "go-user-microservice/internal/pkg/errors"
	"strings"
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
	dbConn := db.GetDBConnection(ctx, r.db)

	operationStory.ExternalID = uuid.New().String()
	operationStory.CreatedAt = time.Now()

	query := ` INSERT INTO operations_stories 
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

	var dbError error
	var result sql.Result
	result, dbError = dbConn.NamedExec(query, operationStory)

	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}

	operationStory.ID = uint64(db.LastInsertID(result))

	return nil
}

func (r *OperationStoryRepository) List(
	ctx context.Context,
	queryFilter *filter.OperationStoryFilter,
) ([]dto.OperationStory, int64, error) {
	dbConn := db.GetDBConnection(ctx, r.db)

	var operationStoriesDto []dto.OperationStory
	innerJoinSelect := " FROM operations_stories os INNER JOIN cards c on c.id = os.card_id "
	query := "SELECT * " + innerJoinSelect
	queryCount := "SELECT COUNT(*)" + innerJoinSelect

	var conditions []string
	params := make(map[string]interface{})

	if queryFilter.OperationType != nil {
		conditions = append(conditions, "os.operation_type_id = :operation_type_id")
		params["operation_type_id"] = queryFilter.OperationType
	}

	query += " WHERE " + strings.Join(conditions, " AND ")
	queryCount += " WHERE " + strings.Join(conditions, " AND ")

	namedQuery, args, e := sqlx.Named(query, params)
	if e != nil {
		return nil, 0, e
	}

	if dbError := dbConn.Select(&operationStoriesDto, namedQuery, args...); dbError != nil {
		return nil, 0, dbError
	}

	var count int64
	namedQueryCount, args, e := sqlx.Named(queryCount, params)
	if e != nil {
		return nil, 0, e
	}

	if dbError := dbConn.Get(&count, namedQueryCount, args...); dbError != nil {
		return nil, 0, dbError
	}

	return operationStoriesDto, count, nil
}
