package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/errors"
)

type Role struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *Role {
	return &Role{
		db: db,
	}
}

func (r *Role) AddUserToRole(ctx context.Context, userID uint64, roleID entities.RoleID) error {
	dbConn := db.GetDBConnection(ctx, r.db)
	query := `INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)`
	if _, e := dbConn.Exec(query, userID, roleID); e != nil {
		return errors.DatabaseError(e)
	}
	return nil
}

func (r *Role) AddUserToRoles(ctx context.Context, userID uint64, roles []entities.RoleID) error {
	dbConn := db.GetDBConnection(ctx, r.db)
	query := `INSERT INTO user_roles (user_id, role_id)`
	for _, role := range roles {
		if _, e := dbConn.Exec(query, userID, role); e != nil {
			return errors.DatabaseError(e)
		}
	}
	return nil
}
