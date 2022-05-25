package providers

import (
	"database/sql"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/database"
	"log"
)

type DatabaseConnectionProvider struct {
	coreConnection *sql.DB
	coreDatabase   database.Database
}

func NewDatabaseConnectionProvider(
	config *config.Config,
) *DatabaseConnectionProvider {
	coreConn, e := sql.Open(config.DatabaseDriver, config.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	dbConnectionWrapper := database.NewDatabase(coreConn)

	return &DatabaseConnectionProvider{
		coreConnection: coreConn,
		coreDatabase:   dbConnectionWrapper,
	}
}

func (c *DatabaseConnectionProvider) GetCoreConnection() *sql.DB {
	return c.coreConnection
}

func (c *DatabaseConnectionProvider) GetCoreDatabaseWrapper() database.Database {
	return c.coreDatabase
}
