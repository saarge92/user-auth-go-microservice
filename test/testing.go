package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/pkg/config"
	providersFunction "go-user-microservice/internal/pkg/domain/providers"
	"go-user-microservice/internal/pkg/domain/servers"
	"go-user-microservice/internal/pkg/providers"
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

func CreateTestServer(
	stripeServiceProvider providersFunction.ProvideFunction,
) (servers.ServerInterface, func(), error) {
	serverTest := NewServerTest(stripeServiceProvider)

	if e := serverTest.InitConfig(); e != nil {
		return nil, nil, e
	}
	var configuration *config.Config

	if e := serverTest.GetDIContainer().Invoke(func(config *config.Config) {
		configuration = config
	}); e != nil {
		return nil, nil, e
	}
	connectionName := fmt.Sprintf(DatabaseName+"_%d", connectionCount)
	txdb.Register(connectionName, "mysql", configuration.CoreDatabaseURL)
	connectionCount++
	configuration.DatabaseDriver = connectionName

	if e := serverTest.InitContainer(); e != nil {
		return nil, nil, e
	}
	var connectionProvider *providers.DatabaseConnectionProvider

	if e := serverTest.GetDIContainer().Invoke(
		func(connProvider *providers.DatabaseConnectionProvider) {
			connectionProvider = connProvider
		}); e != nil {
		return nil, nil, e
	}
	if e := serverTest.InitGRPCServers(); e != nil {
		return nil, nil, e
	}
	coreConn := connectionProvider.GetCoreConnection()
	return serverTest, func() {
		if e := coreConn.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", coreConn.DriverName())
	}, nil
}
