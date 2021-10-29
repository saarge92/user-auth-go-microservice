package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/domain/servers"
	"go-user-microservice/internal/app/providers"
)

const (
	Email                   = "test@test.com"
	Inn              uint64 = 7707083893
	Login            string = "user@mail.ru"
	DatabaseName            = "user-database"
	EmulateUserID           = 1
	EmulateUserLogin        = "test@mail.ru"
	EmulateLoginInn  uint64 = 7842349892
)

var connectionCount = 0

func CreateTestServer() (servers.ServerInterface, func(), error) {
	serverTest := NewServerTest()
	e := serverTest.InitConfig()
	if e != nil {
		return nil, nil, e
	}
	var configuration *config.Config
	e = serverTest.GetDIContainer().Invoke(func(config *config.Config) {
		configuration = config
	})
	if e != nil {
		return nil, nil, e
	}
	connectionName := fmt.Sprintf(DatabaseName+"_%d", connectionCount)
	txdb.Register(connectionName, "mysql", configuration.CoreDatabaseURL)
	connectionCount++
	configuration.DatabaseDriver = connectionName
	e = serverTest.InitContainer()
	if e != nil {
		return nil, nil, e
	}
	var connectionProvider *providers.ConnectionProvider
	e = serverTest.GetDIContainer().Invoke(
		func(connProvider *providers.ConnectionProvider) {
			connectionProvider = connProvider
		})
	if e != nil {
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
