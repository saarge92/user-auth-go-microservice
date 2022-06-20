package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/user/domain"
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/app/wallet/dto"
	walletEntities "go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/dictionary"
	repositoryInterface "go-user-microservice/internal/pkg/domain/repositories"
	sharedEntities "go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   walletDomain.WalletRepository
	userRepository     domain.UserRepository
	currencyRepository repositoryInterface.CurrencyRepository
}

func NewWalletService(
	walletRepository walletDomain.WalletRepository,
	userRepository domain.UserRepository,
	currencyRepository repositoryInterface.CurrencyRepository,
) *WalletService {
	return &WalletService{
		walletRepository:   walletRepository,
		userRepository:     userRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *WalletService) Create(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (wallet *walletEntities.Wallet, e error) {
	userRoleDto, currency, e := s.checkCreateWalletData(ctx, form)
	if e != nil {
		return nil, e
	}

	userData := userRoleDto.User

	if form.IsDefault {
		if resetWalletErr := s.resetDefaultWallet(ctx, userData.ID); resetWalletErr != nil {
			return nil, resetWalletErr
		}
	}

	return s.initWallet(ctx, currency.ID, userData.ID, form.IsDefault)
}

func (s *WalletService) Wallets(ctx context.Context) ([]dto.WalletCurrencyDto, error) {
	var user *entities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*entities.User); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	wallets, e := s.walletRepository.ListByUserID(ctx, user.ID)
	if e != nil {
		return nil, e
	}
	return wallets, nil
}

func (s *WalletService) checkCreateWalletData(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (*userDto.UserRole, *sharedEntities.Currency, error) {
	userRoleDto, e := grpc.GetUserWithRolesFromContext(ctx)
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
	walletUserExist, e := s.walletRepository.Exist(ctx, userRoleDto.User.ID, currency.ID)
	if e != nil {
		return nil, nil, e
	}
	if walletUserExist {
		return nil, nil, status.Error(codes.AlreadyExists, errorlists.UserWalletAlreadyExist)
	}

	return userRoleDto, currency, nil
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
