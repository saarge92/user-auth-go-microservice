package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/pkg/errorlists"
)

func validateExpDay() validation.RuleFunc {
	return func(value interface{}) error {
		if day, ok := value.(uint32); ok {
			if day == 0 {
				return fmt.Errorf(errorlists.MustBeMore, "expire_day", 0)
			}
			if day > 31 {
				return fmt.Errorf(errorlists.MustBeLess, "expire_day", 31)
			}
		}
		return fmt.Errorf(errorlists.ConvertError, "expire_day")
	}
}
