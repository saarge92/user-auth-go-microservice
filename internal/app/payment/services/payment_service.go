package services

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/payment/domain"
	paymentDto "go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	"go-user-microservice/internal/app/payment/form"
	"go-user-microservice/internal/app/payment/transformers"
	userEntities "go-user-microservice/internal/app/user/entities"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/errorlists"
	sharedFilter "go-user-microservice/internal/pkg/filter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentService struct {
	operationStoryRepository domain.OperationStoryRepository
	cardRepository           cardDomain.CardRepository
	walletRepository         walletDomain.WalletRepository
	stripeChargeService      stripe.ChargeService
	coreDB                   *sqlx.DB
}

func NewPaymentService(
	operationStoryRepository domain.OperationStoryRepository,
	walletRepository walletDomain.WalletRepository,
	cardRepository cardDomain.CardRepository,
	stripeChargeService stripe.ChargeService,
	coreDB *sqlx.DB,
) *PaymentService {
	return &PaymentService{
		operationStoryRepository: operationStoryRepository,
		walletRepository:         walletRepository,
		cardRepository:           cardRepository,
		stripeChargeService:      stripeChargeService,
		coreDB:                   coreDB,
	}
}

func (s *PaymentService) Deposit(
	ctx context.Context,
	depositInfo *form.Deposit,
	syncChannel chan<- interface{},
) (operationStory *entities.OperationStory, e error) {
	defer close(syncChannel)
	var user *userEntities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*userEntities.User); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}

	card, e := s.cardRepository.OneByCardAndUserID(ctx, depositInfo.CardExternalId, user.ID)
	if e != nil {
		return nil, e
	}
	walletWithCurrencyDto, e := s.walletRepository.OneByExternalIDAndUserID(ctx, depositInfo.WalletExternalId, user.ID)
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
		CustomerID: user.CustomerProviderID,
	}
	chargeResponse, e := s.stripeChargeService.CardCharge(cardChargeDto)
	if e != nil {
		return nil, e
	}
	operationType := entities.DepositOperationType
	balanceAfter := walletWithCurrencyDto.Balance.Add(amount)
	operationStory = &entities.OperationStory{
		Amount:             amount,
		UserID:             user.ID,
		OperationTypeID:    operationType,
		CardID:             card.ID,
		BalanceBefore:      walletWithCurrencyDto.Balance,
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
	var user *userEntities.User
	var convertOk bool
	if user, convertOk = ctx.Value(dictionary.User).(*userEntities.User); !convertOk {
		return nil, 0, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}

	operationType := transformers.FromGRPCOperationType(request.OperationType)
	paymentFilter := &filter.OperationStoryFilter{
		UserID:        user.ID,
		OperationType: operationType,
		Pagination: sharedFilter.Pagination{
			Page:    request.Pagination.Page,
			PerPage: request.Pagination.PerPage,
		},
	}

	return s.operationStoryRepository.List(ctx, paymentFilter)
}
