package test

import (
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	config2 "go-user-microservice/internal/app/config"
	servers2 "go-user-microservice/internal/app/domain/servers"
	providers2 "go-user-microservice/internal/app/providers"
)

const (
	Inn              uint64 = 7707083893
	Login            string = "member@mail.ru"
	DatabaseName            = "member-database"
	EmulateUserID           = 1
	EmulateUserLogin        = "test@mail.ru"
	EmulateLoginInn  uint64 = 7842349892
)

var connectionCount = 0

func CreateTestServer() (servers2.ServerInterface, func(), error) {
	server := NewServerTest()
	e := server.InitConfig()
	if e != nil {
		return nil, nil, e
	}
	var configuration *config2.Config
	e = server.GetDIContainer().Invoke(func(config *config2.Config) {
		configuration = config
	})
	if e != nil {
		return nil, nil, e
	}
	connectionName := fmt.Sprintf(DatabaseName+"_%d", connectionCount)
	txdb.Register(connectionName, "mysql", configuration.CoreDatabaseURL)
	connectionCount++
	configuration.DatabaseDriver = connectionName
	e = server.InitContainer()
	if e != nil {
		return nil, nil, e
	}
	var connectionProvider *providers2.ConnectionProvider
	e = server.GetDIContainer().Invoke(
		func(connProvider *providers2.ConnectionProvider) {
			connectionProvider = connProvider
		})
	if e != nil {
		return nil, nil, e
	}
	coreConn := connectionProvider.GetCoreConnection()
	return server, func() {
		if e := coreConn.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", coreConn.DriverName())
	}, nil
}
