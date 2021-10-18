package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/test"
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
	query := `INSERT INTO users (
				id, login, inn, name, password, created_at, updated_at)`
	userEntity := &entites.User{
		ID:    test.EmulateUserID,
		Login: test.EmulateUserLogin,
		Inn:   test.EmulateLoginInn,
	}
	_, e := db.NamedExec(query, userEntity)
	if e != nil {
		return e
	}
	return nil
}
