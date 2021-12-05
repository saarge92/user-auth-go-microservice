package domain

import (
	"context"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
)

type PaymentService interface {
	Deposit(ctx context.Context, depositInfo *form.Deposit, syncChannel chan<- interface{}) (*entities.OperationStory, error)
}
