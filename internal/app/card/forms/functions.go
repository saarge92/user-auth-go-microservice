package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/pkg/errorlists"
)

func ValidateExpireMonth() validation.RuleFunc {
	return func(value interface{}) error {
		if month, ok := value.(uint32); ok {
			if month == 0 {
				return fmt.Errorf(errorlists.MustBeMore, "expire_day", 0)
			}
			if month > 12 {
				return fmt.Errorf(errorlists.MustBeLess, "expire_day", 12)
			}
			return nil
		}
		return fmt.Errorf(errorlists.ConvertError, "expire_month")
	}
}
