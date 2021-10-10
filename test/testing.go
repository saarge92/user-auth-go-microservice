package test

import (
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/providers"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	Inn          uint64 = 7707083893
	Login        string = "user@mail.ru"
	DatabaseName        = "user-database"
)

func CreateTestServer() (*app.Server, func()) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	e := os.Chdir(dir)
	if e != nil {
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
		panic(e)
	}
	txdb.Register(DatabaseName, "mysql", configuration.CoreDatabaseURL)
	var connProviderTest *providers.ConnectionProvider
	e = server.GetDIContainer().Invoke(func(connProvider *providers.ConnectionProvider) {
		connProviderTest = connProvider
	})
	coreConn := connProviderTest.GetCoreConnection()
	e = coreConn.Close()
	if e != nil {
		panic(e)
	}
	coreConn, e = sqlx.Open(DatabaseName, configuration.CoreDatabaseURL)
	if e != nil {
		panic(e)
	}
	connProviderTest.SetCoreConnection(coreConn)
	return server, func() {
		time.Sleep(time.Second)
		if e := coreConn.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", coreConn.DriverName())
	}
}
