package domain

import (
	"context"
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
)

type OperationStoryRepository interface {
	Create(ctx context.Context, record *entities.OperationStory) error
	List(
		ctx context.Context,
		queryFilter *filter.OperationStoryFilter) ([]dto.OperationStory, int64, error)
}
