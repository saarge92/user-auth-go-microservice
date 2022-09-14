package domain

import (
	"context"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
	sharedDto "go-user-microservice/internal/pkg/dto"
)

type PaymentService interface {
	Deposit(ctx context.Context, depositInfo form.Deposit) (*entities.OperationStory, error)
	List(ctx context.Context, request *form.ListPayment) ([]dto.OperationStory, int64, error)
}

type StripeChargeService interface {
	CardCharge(cardInfo sharedDto.StripeCardCustomerChargeCreate) (*stripe.Charge, error)
}
