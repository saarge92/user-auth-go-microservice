package payment

import (
	"context"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcServerPayment struct {
	paymentService domain.PaymentService
}

func NewGrpcPaymentServer(
	paymentService domain.PaymentService,
) *GrpcServerPayment {
	return &GrpcServerPayment{
		paymentService: paymentService,
	}
}

func (s *GrpcServerPayment) Deposit(
	ctx context.Context,
	request *core.DepositRequest,
) (*core.DepositResponse, error) {
	depositInfo := &form.Deposit{DepositRequest: request}
	var operationStory *entities.OperationStory
	var depositError error
	syncChannel := make(chan interface{})
	go func() {
		operationStory, depositError = s.paymentService.Deposit(ctx, depositInfo, syncChannel)
	}()
	<-syncChannel
	if depositError != nil {
		return nil, depositError
	}
	return &core.DepositResponse{
		TransactionId: operationStory.ExternalID,
	}, nil
}
