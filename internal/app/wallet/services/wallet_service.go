package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/app/wallet/dto"
	walletEntities "go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/app/wallet/forms"
	sharedEntities "go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   domain.WalletRepository
	currencyRepository domain.CurrencyRepository
}

func NewWalletService(
	walletRepository domain.WalletRepository,
	currencyRepository domain.CurrencyRepository,
) *WalletService {
	return &WalletService{
		walletRepository:   walletRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *WalletService) Create(ctx context.Context, form *forms.WalletCreateForm) (wallet *walletEntities.Wallet, e error) {
	userData, currency, e := s.checkCreateWalletData(ctx, form)
	if e != nil {
		return nil, e
	}

	if form.IsDefault {
		if resetWalletErr := s.resetDefaultWallet(ctx, userData.ID); resetWalletErr != nil {
			return nil, resetWalletErr
		}
	}

	return s.initWallet(ctx, currency.ID, userData.ID, form.IsDefault)
}

func (s *WalletService) MyWallets(ctx context.Context) ([]dto.WalletCurrencyDto, error) {
	userData, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, e
	}

	wallets, e := s.walletRepository.ListByUserID(ctx, userData.ID)
	if e != nil {
		return nil, e
	}
	return wallets, nil
}

func (s *WalletService) checkCreateWalletData(ctx context.Context, form *forms.WalletCreateForm) (*entities.User, *sharedEntities.Currency, error) {
	userData, e := grpc.GetUserWithRolesFromContext(ctx)
	if e != nil {
		return nil, nil, e
	}

	currency, e := s.currencyRepository.GetByCode(ctx, form.Code)
	if e != nil {
		return nil, nil, e
	}
	if currency == nil {
		return nil, nil, status.Error(codes.NotFound, errorlists.CurrencyNotFound)
	}
	walletUserExist, e := s.walletRepository.Exist(ctx, userData.ID, currency.ID)
	if e != nil {
		return nil, nil, e
	}
	if walletUserExist {
		return nil, nil, status.Error(codes.AlreadyExists, errorlists.UserWalletAlreadyExist)
	}

	return userData, currency, nil
}

func (s *WalletService) resetDefaultWallet(ctx context.Context, userID uint64) error {
	defaultWallet, e := s.walletRepository.DefaultWalletByUser(ctx, userID)
	if e != nil {
		return e
	}
	if defaultWallet != nil {
		if e = s.walletRepository.SetAsDefaultForUserWallet(ctx, userID, false); e != nil {
			return e
		}
	}

	return nil
}

func (s *WalletService) initWallet(ctx context.Context, currencyID uint32, userID uint64, isDefault bool) (*walletEntities.Wallet, error) {
	balance := decimal.NewFromInt(0)
	wallet := &walletEntities.Wallet{
		UserID:     userID,
		Balance:    balance,
		CurrencyID: currencyID,
		IsDefault:  isDefault,
	}

	if e := s.walletRepository.Create(ctx, wallet); e != nil {
		return nil, e
	}

	return wallet, nil
}
