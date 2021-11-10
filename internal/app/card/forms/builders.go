package forms

import "go-user-microservice/pkg/protobuf/card"

type CardFormBuilder struct{}

func (b *CardFormBuilder) CreateCreateForm(request *card.CreateCardRequest) *CreateCard {
	expDayValidateRule := validateExpireMonth()
	return NewCreateCardForm(request, expDayValidateRule)
}
