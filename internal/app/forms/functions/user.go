package functions

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/errorlists"
	appErrors "go-user-microservice/internal/app/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserAlReadyExists(userRepository repositories.UserRepositoryInterface) validation.RuleFunc {
	return func(value interface{}) error {
		email := value.(string)
		exist, e := userRepository.UserExist(email)
		if e != nil {
			return appErrors.DatabaseError(e)
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists.UserEmailAlreadyExist)
		}
		return nil
	}
}

func UserWithInnAlreadyExists(userRepository repositories.UserRepositoryInterface) validation.RuleFunc {
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
