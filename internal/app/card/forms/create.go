package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"go-user-microservice/pkg/protobuf/card"
)

type CreateCard struct {
	*card.CreateCardRequest
	expMonthValidateRule validation.RuleFunc
}

func NewCreateCardForm(
	request *card.CreateCardRequest,
	expMonthValidateRule validation.RuleFunc,
) *CreateCard {
	return &CreateCard{
		CreateCardRequest:    request,
		expMonthValidateRule: expMonthValidateRule,
	}
}

func (f *CreateCard) Validate() error {
	return validation.ValidateStruct(
		validation.Field(&f.CardNumber, validation.Required, is.CreditCard),
		validation.Field(&f.ExpireMonth, validation.Required, validation.By(f.expMonthValidateRule)),
		validation.Field(&f.ExpireYear, validation.Required, v.Min(1)),
	)
}
