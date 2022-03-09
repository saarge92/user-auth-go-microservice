package main

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	cardEntities "go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/user/entities"
	walletEntities "go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/test"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func main() {
	if err := godotenv.Load(".env.test"); err != nil {
		panic(err)
	}
	appConfig := config.NewConfig()
	mainDB, e := sqlx.Open("mysql", appConfig.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	if e = createTestUser(mainDB); e != nil {
		log.Fatal(e)
	}
	if e = createTestCreditCard(mainDB); e != nil {
		log.Fatal(e)
	}
	if e = createWallet(mainDB); e != nil {
		log.Fatal(e)
	}
}

func createTestUser(db *sqlx.DB) error {
	now := time.Now()
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(test.UserPasswordForRealUser), bcrypt.DefaultCost)
	if e != nil {
		return e
	}
	query := `INSERT INTO users (
				id, login, inn, name, password, created_at, updated_at)
				VALUES (:id, :login, :inn, :name, :password, :created_at, :updated_at)`
	userEntity := &entities.User{
		ID:        test.UserIDForRealUser,
		Login:     test.UserLoginForRealUser,
		Inn:       test.InnForRealUser,
		Name:      test.NameForRealUser,
		Password:  string(passwordHash),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, e = db.NamedExec(query, userEntity); e != nil {
		return e
	}
	return nil
}

func createTestCreditCard(db *sqlx.DB) error {
	now := time.Now()
	query := `INSERT INTO cards (
			external_id, number, expire_month, expire_year,
			user_id, external_provider_id,
			is_default, created_at, updated_at)
			VALUES (:external_id, :number, :expire_month, :expire_year,
					:user_id, :external_provider_id, :is_default, :created_at, :updated_at)`
	expireYear := uint32(time.Now().Year() + 3)
	externalProviderID := uuid.New().String()
	cardEntity := &cardEntities.Card{
		ExternalID:         test.ExternalIDForCard,
		Number:             test.CardNumberForRealUser,
		ExpireMonth:        test.ExpireMonth,
		ExpireYear:         expireYear,
		UserID:             test.UserIDForRealUser,
		ExternalProviderID: externalProviderID,
		IsDefault:          false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if _, e := db.NamedExec(query, cardEntity); e != nil {
		return e
	}
	return nil
}

func createWallet(db *sqlx.DB) error {
	now := time.Now()
	query := `INSERT INTO wallets(
			external_id, user_id, currency_id, balance,
			is_default, created_at, updated_at)
			VALUES (:external_id, :user_id, :currency_id, :balance, :is_default, 
					:created_at, :updated_at)
			`
	balance := decimal.NewFromInt(1000)
	walletEntity := walletEntities.Wallet{
		ExternalID: test.ExternalIDForWallet,
		CurrencyID: test.USDCurrencyID,
		UserID:     test.UserIDForRealUser,
		Balance:    balance,
		IsDefault:  false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if _, e := db.NamedExec(query, walletEntity); e != nil {
		return e
	}
	return nil
}
