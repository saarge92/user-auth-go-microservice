package forms

import "go-user-microservice/pkg/protobuf/core"

type CardFormBuilder struct{}

func (b *CardFormBuilder) CreateCreateForm(request *core.CreateCardRequest) *CreateCard {
	expDayValidateRule := validateExpireMonth()
	return NewCreateCardForm(request, expDayValidateRule)
}
