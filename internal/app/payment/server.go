package payment

import (
	"context"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/internal/app/wallet/transformer"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcServerPayment struct {
	paymentService     domain.PaymentService
	transactionHandler *db.TransactionHandlerDB
}

func NewGrpcPaymentServer(
	paymentService domain.PaymentService,
	transactionHandler *db.TransactionHandlerDB,
) *GrpcServerPayment {
	return &GrpcServerPayment{
		paymentService:     paymentService,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcServerPayment) Deposit(ctx context.Context, request *core.DepositRequest) (response *core.DepositResponse, e error) {
	ctx, tx, e := s.transactionHandler.Create(ctx, nil)
	if e != nil {
		return nil, e
	}
	defer func() {
		e = db.HandleTransaction(tx, e)
	}()

	depositInfo := form.Deposit{DepositRequest: request}
	operationStory, e := s.paymentService.Deposit(ctx, depositInfo)

	if e != nil {
		return nil, e
	}
	return &core.DepositResponse{
		TransactionId: operationStory.ExternalID,
	}, nil
}

func (s *GrpcServerPayment) List(
	ctx context.Context,
	request *core.ListRequest,
) (*core.ListResponse, error) {
	listRequest := &form.ListPayment{ListRequest: request}
	if e := listRequest.Validate(); e != nil {
		return nil, e
	}

	response, count, e := s.paymentService.List(ctx, listRequest)
	if e != nil {
		return nil, e
	}

	responseGrpc := transformer.FromOperationStoriesDtoToGRPCResponse(
		response,
		count,
	)

	return responseGrpc, nil
}
