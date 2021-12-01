package domain

import (
	"context"
	"go-user-microservice/internal/app/payment/entities"
)

type TransactionRepositoryInterface interface {
	Create(ctx context.Context, record *entities.Transaction) error
}
