package providers

import (
	"database/sql"
	"go-user-microservice/internal/pkg/database"
)

type DatabaseConnectionProvider interface {
	GetCoreConnection() *sql.DB
	GetCoreDatabaseWrapper() database.Database
}
