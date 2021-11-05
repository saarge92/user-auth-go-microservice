package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"go-user-microservice/pkg/protobuf/card"
)

type CreateCard struct {
	*card.CreateCardRequest
	expDayValidateRule   validation.RuleFunc
	expMonthValidateRule validation.RuleFunc
}

func NewCreateCardForm(
	request *card.CreateCardRequest,
	expDayValidateRule validation.RuleFunc,
	expMonthValidateRule validation.RuleFunc,
) *CreateCard {
	return &CreateCard{
		CreateCardRequest:    request,
		expMonthValidateRule: expDayValidateRule,
		expDayValidateRule:   expMonthValidateRule,
	}
}

func (f *CreateCard) Validate() error {
	return validation.ValidateStruct(
		validation.Field(&f.CardNumber, validation.Required, is.CreditCard),
		validation.Field(&f.ExpireDay, validation.Required, validation.By(f.expDayValidateRule)),
		validation.Field(&f.ExpireMonth, validation.Required, validation.By(f.expMonthValidateRule)),
		validation.Field(&f.ExpireYear, validation.Required, v.Min(1)),
	)
}
