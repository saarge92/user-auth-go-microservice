package test

import (
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app"
	"go-user-microservice/internal/config"
	"os"
	"path"
	"runtime"
)

const (
	Inn          uint64 = 7707083893
	Login        string = "user@mail.ru"
	DatabaseName        = "user-database"
)

func CreateTestServer() *app.Server {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	e := os.Chdir(dir)
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
	if e := godotenv.Load(".env.test"); e != nil {
		panic(e)
	}
	server := app.NewServer()
	var configuration *config.Config
	e = server.GetDIContainer().Invoke(func(config *config.Config) {
		configuration = config
	})
	if e != nil {
		log.Fatal(e)
	}
	txdb.Register(DatabaseName, "mysql", configuration.CoreDatabaseURL)
	return server
}
