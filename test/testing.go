package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
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
)

var connectionCount = 0
var databaseConnectionInstance *appProvider.DatabaseConnectionProvider

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
	if databaseConnectionInstance == nil {
		databaseConnectionInstance = appProvider.NewConnectionProvider(appConfig)
	}
	repositoryProvider := appProvider.NewRepositoryProvider(databaseConnectionInstance)
	if stripeServiceProvider == nil {
		stripeServiceProvider = &testProviders.TestStripeServiceProvider{}
	}
	serviceProvider := appProvider.NewServiceProvider(
		appConfig,
		repositoryProvider,
		databaseConnectionInstance,
		stripeServiceProvider,
	)

	coreConn := databaseConnectionInstance.GetCoreConnection()
	grpcServerProvider := appProvider.NewGrpcServerProvider(serviceProvider)

	return grpcServerProvider, func() {
		if e := coreConn.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", coreConn.DriverName())
	}
}
