package repositories

import (
	"context"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/errors"
)

type Role struct {
	databaseInstance database.Database
}

func NewRoleRepository(databaseInstance database.Database) *Role {
	return &Role{
		databaseInstance: databaseInstance,
	}
}

func (r *Role) AddUserToRole(ctx context.Context, userID uint64, roleID entities.RoleID) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)`
	if _, e := r.databaseInstance.ExecContext(ctx, query, userID, roleID); e != nil {
		return errors.DatabaseError(e)
	}
	return nil
}

func (r *Role) AddUserToRoles(ctx context.Context, userID uint64, roles []entities.RoleID) error {
	query := `INSERT INTO user_roles (user_id, role_id)`
	for _, role := range roles {
		if _, e := r.databaseInstance.ExecContext(ctx, query, userID, role); e != nil {
			return errors.DatabaseError(e)
		}
	}
	return nil
}
