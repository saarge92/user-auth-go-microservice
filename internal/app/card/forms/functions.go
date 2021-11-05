package forms

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/pkg/errorlists"
)

func validateExpDay() validation.RuleFunc {
	return func(value interface{}) error {
		if day, ok := value.(uint32); ok {
			if day == 0 {
				return errors.New(fmt.Sprintf(errorlists.MustBeMore, "expire_day", 0))
			}
			if day > 31 {
				return errors.New(fmt.Sprintf(errorlists.MustBeLess, "expire_day", 31))
			}
		}
		return errors.New(fmt.Sprintf(errorlists.ConvertError, "expire_day"))
	}
}
