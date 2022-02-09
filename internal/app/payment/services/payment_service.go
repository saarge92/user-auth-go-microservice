package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
	userEntities "go-user-microservice/internal/app/user/entities"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/repositories"
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
	tx := s.coreDB.MustBegin()
	defer func() {
		e = repositories.HandleTransaction(tx, e)
	}()
	newCtx := context.WithValue(ctx, repositories.CurrentTransaction, tx)
	card, e := s.cardRepository.OneByCardAndUserID(newCtx, depositInfo.CardExternalId, user.ID)
	if e != nil {
		return nil, e
	}
	if card == nil {
		return nil, status.Error(codes.NotFound, errorlists.CardNotFound)
	}
	walletWithCurrencyDto, e := s.walletRepository.OneByExternalIDAndUserID(newCtx, depositInfo.WalletExternalId, user.ID)
	if e != nil {
		return nil, e
	}
	amount := decimal.NewFromFloat(depositInfo.Amount)
	e = s.walletRepository.IncreaseBalanceByID(newCtx, walletWithCurrencyDto.Wallet.ID, amount)
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
	if operationStoryError := s.operationStoryRepository.Create(newCtx, operationStory); operationStoryError != nil {
		return nil, operationStoryError
	}
	return operationStory, nil
}