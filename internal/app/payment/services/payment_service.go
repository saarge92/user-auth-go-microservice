package services

import (
	"context"
	cardDomain "go-user-microservice/internal/app/card/domain"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/form"
	userEntities "go-user-microservice/internal/app/user/entities"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *PaymentService) Deposit(ctx context.Context, depositInfo *form.Deposit) (*entities.OperationStory, error) {
	var user *userEntities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*userEntities.User); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	card, e := s.cardRepository.OneByCardAndUserID(ctx, depositInfo.CardExternalId, user.ID)
	if e != nil {
		return nil, e
	}
	if card == nil {
		return nil, status.Error(codes.NotFound, errorlists.CardNotFound)
	}
	_, e = s.walletRepository.OneByExternalIDAndUserID(ctx, depositInfo.WalletExternalId, user.ID)
	if e != nil {
		return nil, e
	}
	return nil, nil
}
