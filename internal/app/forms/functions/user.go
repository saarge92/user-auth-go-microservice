package functions

import (
	validation "github.com/go-ozzo/ozzo-validation"
	repositories2 "go-user-microservice/internal/app/domain/repositories"
	errorlists2 "go-user-microservice/internal/app/errorlists"
	errors2 "go-user-microservice/internal/app/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserAlReadyExists(userRepository repositories2.UserRepositoryInterface) validation.RuleFunc {
	return func(value interface{}) error {
		email := value.(string)
		exist, e := userRepository.UserExist(email)
		if e != nil {
			return errors2.DatabaseError(e)
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists2.UserEmailAlreadyExist)
		}
		return nil
	}
}

func UserWithInnAlreadyExists(userRepository repositories2.UserRepositoryInterface) validation.RuleFunc {
	return func(value interface{}) error {
		inn := value.(uint64)
		exist, e := userRepository.UserByInnExist(inn)
		if e != nil {
			return e
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists2.UserInnAlreadyExist)
		}
		return nil
	}
}
