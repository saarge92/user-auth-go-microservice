package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/errorlists"
	customErrors "go-user-microservice/internal/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type WalletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
	dbConn := db.GetDBConnection(ctx, r.db)

	wallet.UpdatedAt = time.Now()
	wallet.CreatedAt = time.Now()
	wallet.ExternalID = uuid.New().String()
	query := `INSERT INTO wallets(
                    currency_id, user_id, balance, 
                    external_id, is_default, created_at, updated_at)
				VALUES (:currency_id, :user_id, :balance,
				       :external_id, :is_default, :created_at, :updated_at)`
	var result sql.Result
	var dbError error

	result, dbError = dbConn.NamedExec(query, wallet)
	if dbError != nil {
		return dbError
	}

	wallet.ID = uint64(db.LastInsertID(result))
	return nil
}

func (r *WalletRepository) Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error) {
	query := `SELECT * FROM wallets WHERE user_id = ? AND currency_id = ?`
	wallet := &entities.Wallet{}

	dbConn := db.GetDBConnection(ctx, r.db)

	if dbError := dbConn.Get(wallet, query, userID, currencyID); dbError != nil {
		if dbError == sql.ErrNoRows {
			return false, nil
		}
		return false, dbError
	}
	return true, nil
}

func (r *WalletRepository) ByUserAndDefault(
	ctx context.Context,
	userID uint64,
	isDefault bool,
) (*entities.Wallet, error) {
	dbConn := db.GetDBConnection(ctx, r.db)

	query := `SELECT * FROM wallets WHERE user_id = ? AND is_default = ?`
	wallet := &entities.Wallet{}

	if dbError := dbConn.Get(wallet, query, userID, isDefault); dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, status.New(codes.NotFound, errorlists.WalletNotFound).Err()
		}
		return nil, customErrors.DatabaseError(dbError)
	}
	return wallet, nil
}

func (r *WalletRepository) UpdateStatusByUserID(ctx context.Context, userID uint64, isDefault bool) error {
	query := `UPDATE wallets SET is_default = ? WHERE user_id = ?`
	dbConn := db.GetDBConnection(ctx, r.db)

	if _, dbError := dbConn.Exec(query, isDefault, userID); dbError != nil {
		return dbError
	}
	return nil
}

func (r *WalletRepository) ListByUserID(ctx context.Context, userID uint64) ([]dto.WalletCurrencyDto, error) {
	dbConn := db.GetDBConnection(ctx, r.db)
	query := `SELECT 
				wallets.id "wallet.id",
       			wallets.external_id "wallet.external_id",
       			wallets.balance "wallet.balance",
       			wallets.is_default "wallet.balance",
       			currencies.code "currency.code"
				FROM wallets 
				INNER JOIN currencies ON wallets.currency_id = currencies.id
				WHERE wallets.user_id = ?
			`
	var walletsCurrencies []dto.WalletCurrencyDto

	if dbError := dbConn.Select(walletsCurrencies, query, userID); dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, nil
		}
		return nil, dbError
	}
	return walletsCurrencies, nil
}

func (r *WalletRepository) OneByExternalIDAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*dto.WalletCurrencyDto, error) {
	dbConn := db.GetDBConnection(ctx, r.db)
	query := `SELECT 
				wallets.id "wallet.id",
       			wallets.external_id "wallet.external_id",
       			wallets.balance "wallet.balance",
       			wallets.is_default "wallet.balance",
       			currencies.code "currency.code"
				FROM wallets 
				INNER JOIN currencies ON wallets.currency_id = currencies.id
				WHERE wallets.user_id = ? AND wallets.external_id = ?
			`
	walletWithCurrencyDto := new(dto.WalletCurrencyDto)

	if dbError := dbConn.Get(walletWithCurrencyDto, query, userID, externalID); dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, customErrors.CustomDatabaseError(codes.NotFound, errorlists.WalletNotFound)
		}
		return nil, customErrors.DatabaseError(dbError)
	}

	return walletWithCurrencyDto, nil
}

func (r *WalletRepository) IncreaseBalanceByID(ctx context.Context, id uint64, amount decimal.Decimal) error {
	dbConn := db.GetDBConnection(ctx, r.db)

	updatedAt := time.Now()
	query := `UPDATE wallets SET balance = balance + ?,
				updated_at = ? WHERE id = ?`
	var dbError error
	var result sql.Result
	result, dbError = dbConn.Exec(query, amount, updatedAt, id)

	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return customErrors.CustomDatabaseError(codes.NotFound, errorlists.WalletNotFound)
	}

	return nil
}
