package test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/database"
	"os"
	"path"
	"runtime"
)

const (
	InnForTest              uint64 = 7707083893
	LoginForTest            string = "user@mail.ru"
	DatabaseName                   = "user-database"
	UserIDForRealUser       uint64 = 1
	UserLoginForRealUser           = "test@mail.ru"
	UserPasswordForRealUser        = "test123!`"
	InnForRealUser          uint64 = 7842349892
	NameForRealUser                = "Ivan"
	CurrencyCode                   = "RUB"
	CardNumber                     = "4000056655665556"
	ExpireMonth                    = 03
	CVC                            = 366
	CardNumberForRealUser          = "5200828282828210"
	ExternalIDForCard              = "949fcc24-1f5e-45a3-a057-770c04475eb8"
	ExternalIDForWallet            = "fe65292c-3d24-4536-97b7-57f4d3ba63b5"
	USDCurrencyID                  = 2
)

func LoadTestEnv() error {
	// Setting root directory for tests as project root
	// For more details read this topic https://brandur.org/fragments/testing-go-project-root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	if err = godotenv.Load(".env.test"); err != nil {
		return err
	}
	return nil
}

func InitConnectionsWithCloseFunc() (*sql.DB, func()) {
	appConfig := config.NewConfig()
	dbConn, txDBCloseFunc := database.TestDBConnection(appConfig.CoreDatabaseURL)

	return dbConn, txDBCloseFunc
}
