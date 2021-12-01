package services

import (
	"context"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
)

type PaymentService struct {
}

func (s *PaymentService) Deposit(ctx context.Context, f *form.Deposit) (*entities.Transaction, error) {
	return nil, nil
}
