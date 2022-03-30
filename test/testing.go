package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/db"
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
	CardNumberForRealUser          = "5200828282828210"
	ExternalIDForCard              = "949fcc24-1f5e-45a3-a057-770c04475eb8"
	ExternalIDForWallet            = "fe65292c-3d24-4536-97b7-57f4d3ba63b5"
	USDCurrencyID                  = 2
)

var connectionCount = 0

func CreateTestServer(
	stripeServiceProvider providers.StripeServiceProvider,
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
	dbConnProvider := appProvider.NewDatabaseConnectionProvider(appConfig)
	mainConnection := dbConnProvider.GetCoreConnection()
	mainConnection.SetMaxIdleConns(1)
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

	transactionHandler := db.NewTransactionHandler(dbConnProvider.GetCoreConnection())
	grpcServerProvider := appProvider.NewGrpcServerProvider(serviceProvider, transactionHandler)

	return grpcServerProvider, func() {
		if e := mainConnection.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", mainConnection.DriverName())
	}
}
