package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/domain/providers"
	appProvider "go-user-microservice/internal/pkg/providers"
	testProviders "go-user-microservice/test/mocks/providers"
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
)

var connectionCount = 0

func CreateTestServer(
	stripeServiceProvider providers.StripeServiceProviderInterface,
) (*appProvider.GrpcServerProvider, func()) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	if err = godotenv.Load(".env.test"); err != nil {
		panic(err)
	}

	appConfig := config.NewConfig()

	connectionName := fmt.Sprintf(DatabaseName+"_%d", connectionCount)
	txdb.Register(connectionName, "mysql", appConfig.CoreDatabaseURL)
	connectionCount++
	appConfig.DatabaseDriver = connectionName
	db, err := sqlx.Open(connectionName, appConfig.CoreDatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(1)
	dbConnProvider := appProvider.NewDatabaseConnectionProvider(appConfig)
	repositoryProvider := appProvider.NewRepositoryProvider(dbConnProvider)
	if stripeServiceProvider == nil {
		stripeServiceProvider = &testProviders.TestStripeServiceProvider{}
	}
	serviceProvider := appProvider.NewServiceProvider(
		appConfig,
		repositoryProvider,
		dbConnProvider,
		stripeServiceProvider,
	)

	grpcServerProvider := appProvider.NewGrpcServerProvider(serviceProvider)

	return grpcServerProvider, func() {
		if e := db.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", db.DriverName())
	}
}
