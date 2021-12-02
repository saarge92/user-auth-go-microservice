package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/internal/app/wallet/entities"
	customErrors "go-user-microservice/internal/pkg/errors"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
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
	tx := sharedRepositories.GetDBTransaction(ctx)
	if tx != nil {
		result, dbError = tx.NamedExec(query, wallet)
	} else {
		result, dbError = r.db.NamedExec(query, wallet)
	}
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	wallet.ID = uint64(sharedRepositories.LastInsertID(result))
	return nil
}

func (r *WalletRepository) Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error) {
	query := `SELECT * FROM wallets WHERE user_id = ? AND currency_id = ?`
	wallet := &entities.Wallet{}

	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Get(wallet, query, userID, currencyID)
	} else {
		dbError = r.db.Get(wallet, query, userID, currencyID)
	}
	if dbError != nil {
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
	query := `SELECT * FROM wallets WHERE user_id = ? AND is_default = ?`
	wallet := &entities.Wallet{}
	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Get(wallet, query, userID, isDefault)
	} else {
		dbError = r.db.Get(wallet, query, userID, isDefault)
	}
	if dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, nil
		}
		return nil, customErrors.DatabaseError(dbError)
	}
	return wallet, nil
}

func (r *WalletRepository) UpdateStatusByUserID(ctx context.Context, userID uint64, isDefault bool) error {
	query := `UPDATE wallets SET is_default = ? WHERE user_id = ?`
	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		_, dbError = tx.Exec(query, isDefault, userID)
	} else {
		_, dbError = r.db.Exec(query, isDefault, userID)
	}
	if dbError != nil {
		return dbError
	}
	return nil
}

func (r *WalletRepository) ListByUserID(ctx context.Context, userID uint64) ([]dto.WalletCurrencyDto, error) {
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
	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	var result *sqlx.Rows
	if tx != nil {
		result, dbError = tx.Queryx(query, userID)
	} else {
		result, dbError = r.db.Queryx(query, userID)
	}
	defer result.Close()
	if dbError != nil {
		return nil, customErrors.DatabaseError(dbError)
	}
	for result.Next() {
		var walletCurrencyDto dto.WalletCurrencyDto
		if e := result.StructScan(&walletCurrencyDto); e != nil {
			return nil, customErrors.DatabaseError(e)
		}
		walletsCurrencies = append(walletsCurrencies, walletCurrencyDto)
	}

	return walletsCurrencies, nil
}

func (r *WalletRepository) OneByExternalIDAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*dto.WalletCurrencyDto, error) {
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
	var dbError error
	tx := sharedRepositories.GetDBTransaction(ctx)
	if tx != nil {
		dbError = tx.Get(walletWithCurrencyDto, query, userID, externalID)
	} else {
		dbError = r.db.Get(walletWithCurrencyDto, query, userID, externalID)
	}
	if dbError != nil {
		return nil, customErrors.DatabaseError(dbError)
	}
	return walletWithCurrencyDto, nil
}

func (r *WalletRepository) IncreaseBalanceByID(ctx context.Context, id uint64, amount decimal.Decimal) error {
	updatedAt := time.Now()
	query := `UPDATE wallets SET balance = balance + ?,
				updated_at = ? WHERE id = ?`
	var dbError error
	tx := sharedRepositories.GetDBTransaction(ctx)
	if tx != nil {
		_, dbError = tx.Exec(query, amount, updatedAt, id)
	} else {
		_, dbError = r.db.Exec(query, amount, updatedAt, id)
	}
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	return nil
}
