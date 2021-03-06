package services

import (
	"context"
	"github.com/shopspring/decimal"
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/payment/domain"
	paymentDto "go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/internal/app/payment/transformers"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	sharedFilter "go-user-microservice/internal/pkg/filter"
	"go-user-microservice/internal/pkg/grpc"
)

type PaymentService struct {
	operationStoryRepository domain.OperationStoryRepository
	cardRepository           cardDomain.CardRepository
	walletRepository         walletDomain.WalletRepository
	stripeChargeService      stripe.ChargeService
}

func NewPaymentService(
	operationStoryRepository domain.OperationStoryRepository,
	walletRepository walletDomain.WalletRepository,
	cardRepository cardDomain.CardRepository,
	stripeChargeService stripe.ChargeService,
) *PaymentService {
	return &PaymentService{
		operationStoryRepository: operationStoryRepository,
		walletRepository:         walletRepository,
		cardRepository:           cardRepository,
		stripeChargeService:      stripeChargeService,
	}
}

func (s *PaymentService) Deposit(
	ctx context.Context,
	depositInfo *form.Deposit,
	syncChannel chan<- interface{},
) (operationStory *entities.OperationStory, e error) {
	defer close(syncChannel)

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
	cardChargeDto := &dto.StripeCardCustomerChargeCreate{
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

func (s *PaymentService) List(
	ctx context.Context,
	request *form.ListPayment,
) (response []paymentDto.OperationStory, count int64, e error) {
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
