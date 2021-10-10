package functions

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/errorlists"
	"go-user-microservice/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserAlReadyExists(userRepository repositories.UserRepository) validation.RuleFunc {
	return func(value interface{}) error {
		email := value.(string)
		exist, e := userRepository.UserExist(email)
		if e != nil {
			return errors.DatabaseError(e)
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists.UserEmailAlreadyExist)
		}
		return nil
	}
}

func UserWithInnAlreadyExists(userRepository repositories.UserRepository) validation.RuleFunc {
	return func(value interface{}) error {
		inn := value.(uint64)
		exist, e := userRepository.UserByInnExist(inn)
		if e != nil {
			return e
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists.UserInnAlreadyExist)
		}
		return nil
	}
}
