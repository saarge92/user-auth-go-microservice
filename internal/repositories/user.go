package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/entites"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entites.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	query := `INSERT INTO users (name, login, password, created_at, updated_at)
				VALUES (:name, :login, :password, :created_at, :updated_at)`
	result, e := r.db.NamedExec(query, user)
	if e != nil {
		return e
	}
	user.ID = lastInsertID(result)
	return nil
}

func (r *UserRepository) Update(user *entites.User) error {
	now := time.Now()
	user.UpdatedAt = now
	query := `UPDATE users SET
				login = :login, password = :password, name = :name,
				created_at = :created_at, updated_at = :updated_at, is_banned = :is_banned`
	_, e := r.db.NamedExec(query, user)
	if e != nil {
		return e
	}
	return nil
}

func (r *UserRepository) UserExist(login string) (bool, error) {
	query := `SELECT * from users where users.login = ?`
	var user = &entites.User{}
	e := r.db.Get(user, query, login)
	if e != nil {
		if e == sql.ErrNoRows {
			return false, nil
		}
		return false, e
	}
	return true, nil
}
