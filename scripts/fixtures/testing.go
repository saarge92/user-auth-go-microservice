package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/user/entities"
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
	e = createTestUser(mainDB)
	if e != nil {
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
	_, e = db.NamedExec(query, userEntity)
	if e != nil {
		return e
	}
	return nil
}
