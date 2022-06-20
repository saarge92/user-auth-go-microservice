package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Knetic/go-namedParameterQuery"
	"github.com/blockloop/scan"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/errorlists"
	sharedErrors "go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/repositories"
	"google.golang.org/grpc/codes"
	"time"
)

type UserRepository struct {
	databaseConnection database.Database
}

func NewUserRepository(db database.Database) *UserRepository {
	return &UserRepository{databaseConnection: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	query := `INSERT INTO users (name, login, inn, password, account_provider_id,
    								customer_provider_id, created_at, updated_at)
				VALUES (:name, :login, :inn, :password, :accountProviderId,
				        	:customerProviderId, :createdAt, :updatedAt)`
	queryNamed := namedParameterQuery.NewNamedParameterQuery(query)
	insertParams := map[string]interface{}{
		"name":               user.Name,
		"login":              user.Login,
		"inn":                user.Inn,
		"password":           user.Password,
		"accountProviderId":  user.AccountProviderID,
		"customerProviderId": user.CustomerProviderID,
		"createdAt":          user.CreatedAt,
		"updatedAt":          user.UpdatedAt,
	}
	queryNamed.SetValuesFromMap(insertParams)
	result, e := r.databaseConnection.ExecContext(ctx, queryNamed.GetParsedQuery(), queryNamed.GetParsedParameters()...)
	if e != nil {
		return e
	}
	user.ID = uint64(repositories.LastInsertID(result))
	return nil
}

func (r *UserRepository) UserExist(ctx context.Context, login string) (bool, error) {
	query := `SELECT * from users where users.login = ?`
	var user = &entities.User{}
	result, e := r.databaseConnection.QueryContext(ctx, query, login)
	if e != nil {
		return false, sharedErrors.DatabaseError(e)
	}

	if e = scan.Row(user, result); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return false, nil
		}
		return false, sharedErrors.DatabaseError(e)
	}

	return true, nil
}

func (r *UserRepository) GetUserWithRoles(ctx context.Context, login string) (*dto.UserRole, error) {
	queryUserSelect := `SELECT * FROM users where login = ?`
	user := new(entities.User)
	userRow, userError := r.databaseConnection.QueryContext(ctx, queryUserSelect, login)
	if userError != nil {
		return nil, sharedErrors.DatabaseError(userError)
	}

	if e := scan.Row(user, userRow); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, sharedErrors.CustomDatabaseError(codes.NotFound, errorlists.UserNotFound)
		}
		return nil, sharedErrors.DatabaseError(userError)
	}

	queryRolesSelect := `SELECT * FROM roles INNER JOIN user_roles
							ON user_roles.role_id = roles.id AND user_roles.user_id = ?`
	var roles []entities.Role
	roleRows, roleError := r.databaseConnection.QueryContext(ctx, queryRolesSelect, user.ID)

	if roleError != nil {
		return nil, sharedErrors.DatabaseError(roleError)
	}

	if e := scan.Rows(&roles, roleRows); e != nil {
		if !errors.Is(e, sql.ErrNoRows) {
			return nil, sharedErrors.DatabaseError(e)
		}
	}

	return &dto.UserRole{
		User:  *user,
		Roles: roles,
	}, nil
}

func (r *UserRepository) UserByInnOrLoginExist(ctx context.Context, login string, inn uint64) (bool, error) {
	query := `SELECT COUNT(*) > 0 FROM users WHERE login = ? OR inn = ?`
	var exist bool

	if e := r.databaseConnection.QueryRowContext(ctx, query, login, inn).Scan(&exist); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return false, nil
		}
		return false, sharedErrors.DatabaseError(e)
	}

	return exist, nil
}
