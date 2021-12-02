package payment

import (
	"context"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/pkg/protobuf/payment"
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
	request *payment.DepositRequest,
) (*payment.DepositResponse, error) {
	depositInfo := &form.Deposit{DepositRequest: request}
	_, e := s.paymentService.Deposit(ctx, depositInfo)
	if e != nil {
		return nil, e
	}
	return nil, nil
}
