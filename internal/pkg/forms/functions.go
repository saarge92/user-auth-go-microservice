package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/pkg/errors"
	"go-user-microservice/pkg/protobuf/core"
)

func ValidatePagination(maxPerPage uint32) validation.RuleFunc {
	return func(value interface{}) error {
		if val := value.(*core.Pagination); val != nil {
			return validatePaginationPages(val.Page, val.PerPage, maxPerPage)
		}
		return nil
	}
}

func validatePaginationPages(page uint32, perPage uint32, maxPerPage uint32) error {
	if page == 0 {
		return fmt.Errorf("pagination.page "+errors.MustBeNoLessThan, 0)
	}
	if perPage == 0 {
		return fmt.Errorf("pagination.per_page "+errors.MustBeNoLessThan, 0)
	}
	if perPage > maxPerPage {
		return fmt.Errorf("pagination.per_page "+errors.MustBeNoGreaterThan, maxPerPage)
	}
	return nil
}
