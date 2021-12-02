package domain

import (
	"context"
	"go-user-microservice/internal/app/payment/entities"
)

type OperationStoryRepository interface {
	Create(ctx context.Context, record *entities.OperationStory) error
}
