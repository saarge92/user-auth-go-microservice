package test

import (
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/contracts/servers"
	"go-user-microservice/internal/providers"
	"time"
)

const (
	Inn          uint64 = 7707083893
	Login        string = "user@mail.ru"
	DatabaseName        = "user-database"
)

func CreateTestServer() (servers.ServerInterface, func(), error) {
	server := NewServerTest()
	e := server.InitConfig()
	if e != nil {
		return nil, nil, e
	}
	var configuration *config.Config
	e = server.GetDIContainer().Invoke(func(config *config.Config) {
		configuration = config
	})
	if e != nil {
		return nil, nil, e
	}
	txdb.Register(DatabaseName, "mysql", configuration.CoreDatabaseURL)
	e = server.InitContainer()
	if e != nil {
		return nil, nil, e
	}
	var connectionProvider *providers.ConnectionProvider
	e = server.GetDIContainer().Invoke(
		func(connProvider *providers.ConnectionProvider) {
			connectionProvider = connProvider
		})
	if e != nil {
		return nil, nil, e
	}
	coreConn := connectionProvider.GetCoreConnection()
	return server, func() {
		time.Sleep(time.Second)
		if e := coreConn.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", coreConn.DriverName())
	}, nil
}
