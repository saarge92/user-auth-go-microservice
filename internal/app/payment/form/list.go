package form

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/pkg/forms"
	"go-user-microservice/pkg/protobuf/core"
)

const maxPerPage uint32 = 100

type ListPayment struct {
	*core.ListRequest
}

func (f *ListPayment) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Pagination, validation.Required, validation.By(forms.ValidatePagination(maxPerPage))),
		validation.Field(&f.OperationType, validation.In(
			core.OperationType_ALL,
			core.OperationType_Deposit,
			core.OperationType_WithDraw,
			core.OperationType_Refund,
		)),
	)
}
