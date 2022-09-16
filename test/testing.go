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

const ()

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
