package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Knetic/go-namedParameterQuery"
	"github.com/blockloop/scan"
	"github.com/google/uuid"
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	"go-user-microservice/internal/pkg/database"
	"strings"
	"time"
)

type OperationStoryRepository struct {
	databaseInstance database.Database
}

func NewOperationStoryRepository(databaseInstance database.Database) *OperationStoryRepository {
	return &OperationStoryRepository{
		databaseInstance: databaseInstance,
	}
}

func (r *OperationStoryRepository) Create(ctx context.Context, operationStory *entities.OperationStory) error {
	operationStory.ExternalID = uuid.New().String()
	operationStory.CreatedAt = time.Now()

	query := ` INSERT INTO operations_stories 
     			(
     			 	external_id, user_id, card_id, amount, balance_before, balance_after,
					external_provider_id, operation_type_id, created_at
				)
				VALUES 
				(
					:externalId, :userId, :cardId, :amount, :balanceBefore, :balanceAfter,
					:externalProviderId, :operationTypeId, :createdAt
				)
				`
	queryNamed := namedParameterQuery.NewNamedParameterQuery(query)
	insertParams := map[string]interface{}{
		"externalId":         operationStory.ExternalID,
		"userId":             operationStory.UserID,
		"cardId":             operationStory.CardID,
		"amount":             operationStory.Amount,
		"balanceBefore":      operationStory.BalanceBefore,
		"balanceAfter":       operationStory.BalanceAfter,
		"externalProviderId": operationStory.ExternalProviderID,
		"operationTypeId":    operationStory.OperationTypeID,
		"createdAt":          operationStory.CreatedAt,
	}
	queryNamed.SetValuesFromMap(insertParams)

	result, dbError := r.databaseInstance.ExecContext(ctx, queryNamed.GetParsedQuery(), queryNamed.GetParsedParameters()...)

	if dbError != nil {
		return dbError
	}

	operationStory.ID = uint64(database.LastInsertID(result))

	return nil
}

func (r *OperationStoryRepository) List(
	ctx context.Context,
	queryFilter *filter.OperationStoryFilter,
) ([]dto.OperationStory, int64, error) {
	var operationStoriesDto []dto.OperationStory
	innerJoinSelect := " FROM operations_stories os INNER JOIN cards c on c.id = os.card_id "
	query := "SELECT * " + innerJoinSelect
	queryCount := "SELECT COUNT(*)" + innerJoinSelect

	var conditions []string
	params := make(map[string]interface{})

	if queryFilter.OperationType != nil {
		conditions = append(conditions, "os.operation_type_id = :operationTypeId")
		params["operationTypeId"] = queryFilter.OperationType
	}

	if conditions != nil {
		query += " WHERE " + strings.Join(conditions, " AND ")
		queryCount += " WHERE " + strings.Join(conditions, " AND ")
	}

	query = database.AddPagination(query, queryFilter.Pagination)

	queryNamed := namedParameterQuery.NewNamedParameterQuery(query)
	queryNamed.SetValuesFromMap(params)
	listRows, listErr := r.databaseInstance.QueryContext(ctx, queryNamed.GetParsedQuery(), queryNamed.GetParsedParameters()...)
	if listErr != nil {
		return nil, 0, listErr
	}
	if e := scan.Rows(operationStoriesDto, listRows); e != nil {
		if !errors.Is(e, sql.ErrNoRows) {
			return nil, 0, e
		}
	}

	var count int64
	queryCountNamed := namedParameterQuery.NewNamedParameterQuery(queryCount)
	queryCountNamed.SetValuesFromMap(params)

	if dbError := r.databaseInstance.QueryRowContext(
		ctx,
		queryCountNamed.GetParsedQuery(),
		queryCountNamed.GetParsedParameters()...,
	).Scan(&count); dbError != nil {
		return nil, 0, dbError
	}

	return operationStoriesDto, count, nil
}
