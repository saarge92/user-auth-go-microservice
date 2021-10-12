package repositories

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errors"
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

func (r *WalletRepository) Create(wallet *entites.Wallet) error {
	wallet.UpdatedAt = time.Now()
	wallet.CreatedAt = time.Now()
	query := `INSERT INTO wallets(currency_cd, user_id, balance, created_at, updated_at)
				VALUES (:currency_cd, :user_id, :balance, :created_at, :updated_at)`
	result, e := r.db.NamedExec(query, wallet)
	if e != nil {
		return errors.DatabaseError(e)
	}
	wallet.ID = uint64(lastInsertID(result))
	return nil
}
