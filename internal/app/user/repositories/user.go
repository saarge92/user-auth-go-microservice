package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/user/entities"
	sharedErrors "go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/repositories"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	query := `INSERT INTO users (name, login, inn, password, payment_provider_account_id,
    								created_at, updated_at)
				VALUES (:name, :login, :inn, :password, :payment_provider_account_id,
				        	:created_at, :updated_at)`
	result, e := r.db.NamedExec(query, user)
	if e != nil {
		return e
	}
	user.ID = uint64(repositories.LastInsertID(result))
	return nil
}

func (r *UserRepository) Update(user *entities.User) error {
	now := time.Now()
	user.UpdatedAt = now
	query := `UPDATE users SET
				login = :login, password = :password, name = :name, inn = :inn,
				created_at = :created_at, updated_at = :updated_at, is_banned = :is_banned`
	_, e := r.db.NamedExec(query, user)
	if e != nil {
		return e
	}
	return nil
}

func (r *UserRepository) UserExist(login string) (bool, error) {
	query := `SELECT * from users where users.login = ?`
	var user = &entities.User{}
	e := r.db.Get(user, query, login)
	if e != nil {
		if e == sql.ErrNoRows {
			return false, nil
		}
		return false, sharedErrors.DatabaseError(e)
	}
	return true, nil
}

func (r *UserRepository) GetUser(login string) (*entities.User, error) {
	query := `SELECT * FROM users WHERE users.login = ?`
	var user = &entities.User{}
	if e := r.db.Get(user, query, login); e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, sharedErrors.DatabaseError(e)
	}
	return user, nil
}

func (r *UserRepository) UserByInnExist(inn uint64) (bool, error) {
	query := `SELECT * FROM users WHERE users.inn = ?`
	user := &entities.User{}
	if e := r.db.Get(user, query, inn); e != nil {
		if e == sql.ErrNoRows {
			return false, nil
		}
		return false, sharedErrors.DatabaseError(e)
	}
	return true, nil
}

func (r *UserRepository) UserByID(id uint64) (*entities.User, error) {
	query := `SELECT * FROM users WHERE users.id = ?`
	user := &entities.User{}
	if e := r.db.Get(user, query, id); e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, sharedErrors.DatabaseError(e)
	}
	return user, nil
}
