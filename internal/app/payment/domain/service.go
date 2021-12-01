package domain

import (
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
)

type PaymentServiceInterface interface {
	Deposit(depositInfo *form.Deposit) (*entities.Transaction, error)
}
