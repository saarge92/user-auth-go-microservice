package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/pkg/protobuf/core"
)

type SignIn struct {
	*core.SignInMessage
}

func (f *SignIn) Validate() error {
	return validation.ValidateStruct(
		&f,
		validation.Field(&f.Login, validation.Required),
		validation.Field(&f.Password, validation.Required),
	)
}
