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
			return errors.CustomDatabaseError(e)
		}
		if exist {
			return status.Error(codes.AlreadyExists, errorlists.RemoteServerBadAuthorization)
		}
		return nil
	}
}
