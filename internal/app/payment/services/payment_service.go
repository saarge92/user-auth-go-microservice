package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/payment/domain"
	paymentDto "go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/internal/app/payment/transformers"
	"go-user-microservice/internal/pkg/dto"
	sharedFilter "go-user-microservice/internal/pkg/filter"
	"go-user-microservice/internal/pkg/grpc"
)

type Payment struct {
	operationStoryRepository domain.OperationStoryRepository
	cardRepository           domain.CardRepository
	walletRepository         domain.WalletRepository
	stripeChargeService      domain.StripeChargeService
}

func NewPaymentService(
	operationStoryRepository domain.OperationStoryRepository,
	walletRepository domain.WalletRepository,
	cardRepository domain.CardRepository,
	stripeChargeService domain.StripeChargeService,
) *Payment {
	return &Payment{
		operationStoryRepository: operationStoryRepository,
		walletRepository:         walletRepository,
		cardRepository:           cardRepository,
		stripeChargeService:      stripeChargeService,
	}
}

func (s *Payment) Deposit(ctx context.Context, depositInfo form.Deposit) (operationStory *entities.OperationStory, e error) {
	userRoleDto, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}

	card, e := s.cardRepository.OneByCardAndUserID(ctx, depositInfo.CardExternalId, userRoleDto.User.ID)
	if e != nil {
		return nil, e
	}
	walletWithCurrencyDto, e := s.walletRepository.OneByExternalIDAndUserID(ctx, depositInfo.WalletExternalId, userRoleDto.User.ID)
	if e != nil {
		return nil, e
	}
	amount := decimal.NewFromFloat(depositInfo.Amount)
	if e = s.walletRepository.IncreaseBalanceByID(ctx, walletWithCurrencyDto.Wallet.ID, amount); e != nil {
		return nil, e
	}
	cardChargeDto := dto.StripeCardCustomerChargeCreate{
		Amount:     amount,
		Currency:   walletWithCurrencyDto.Currency.Code,
		CardID:     card.ExternalProviderID,
		CustomerID: userRoleDto.User.CustomerProviderID,
	}
	chargeResponse, e := s.stripeChargeService.CardCharge(cardChargeDto)
	if e != nil {
		return nil, e
	}
	operationType := entities.DepositOperationType
	walletInstance := walletWithCurrencyDto.Wallet
	balanceAfter := walletInstance.Balance.Add(amount)
	operationStory = &entities.OperationStory{
		Amount:             amount,
		UserID:             userRoleDto.User.ID,
		OperationTypeID:    operationType,
		CardID:             card.ID,
		BalanceBefore:      walletInstance.Balance,
		BalanceAfter:       balanceAfter,
		ExternalProviderID: chargeResponse.ID,
	}
	if operationStoryError := s.operationStoryRepository.Create(ctx, operationStory); operationStoryError != nil {
		return nil, operationStoryError
	}
	return operationStory, nil
}

func (s *Payment) List(ctx context.Context, request *form.ListPayment) (response []paymentDto.OperationStory, count int64, e error) {
	user, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, 0, e
	}

	operationType := transformers.FromGRPCOperationType(request.OperationType)
	paymentFilter := &filter.OperationStoryFilter{
		UserID:        user.User.ID,
		OperationType: operationType,
		Pagination: sharedFilter.Pagination{
			Page:    request.Pagination.Page,
			PerPage: request.Pagination.PerPage,
		},
	}

	return s.operationStoryRepository.List(ctx, paymentFilter)
}
