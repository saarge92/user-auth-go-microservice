package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go-user-microservice/pkg/protobuf/core"
)

type CreateCard struct {
	*core.CreateCardRequest
}

func (f *CreateCard) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.CardNumber, validation.Required, is.CreditCard),
		validation.Field(&f.ExpireMonth, validation.Required, validation.By(ValidateExpireMonth())),
		validation.Field(&f.ExpireYear, validation.Required),
		validation.Field(&f.Cvc, validation.Required),
	)
}
