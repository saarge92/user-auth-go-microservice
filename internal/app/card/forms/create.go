package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go-user-microservice/pkg/protobuf/core"
)

type CreateCard struct {
	*core.CreateCardRequest
	expMonthValidateRule validation.RuleFunc
}

func NewCreateCardForm(
	request *core.CreateCardRequest,
	expMonthValidateRule validation.RuleFunc,
) *CreateCard {
	return &CreateCard{
		CreateCardRequest:    request,
		expMonthValidateRule: expMonthValidateRule,
	}
}

func (f *CreateCard) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.CardNumber, validation.Required, is.CreditCard),
		validation.Field(&f.ExpireMonth, validation.Required, validation.By(f.expMonthValidateRule)),
		validation.Field(&f.ExpireYear, validation.Required),
		validation.Field(&f.Cvc, validation.Required),
	)
}
