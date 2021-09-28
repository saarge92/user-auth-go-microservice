package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/entites"
	"go-user-microservice/pkg/protobuf/user"
	"regexp"
)

type SignUp struct {
	*user.SignUpMessage
	userExistRule validation.RuleFunc
}

func NewSignUpForm(
	request *user.SignUpMessage,
	userExistRule validation.RuleFunc) *SignUp {
	return &SignUp{
		request,
		userExistRule,
	}
}

func (f *SignUp) Validate() error {
	//emailPattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
	//	"(?:\\\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	innPattern := fmt.Sprintf(`\S{%d}$`, entites.InnLength)
	namePattern := fmt.Sprintf(`\S{%d,%d}$`, entites.NameLengthMin, entites.NameLengthMax)
	return validation.ValidateStruct(f,
		validation.Field(
			&f.Login,
			validation.Required,
			validation.By(f.userExistRule)),
		validation.Field(
			&f.Password,
			validation.Required,
			validation.Length(6, 120),
		),
		validation.Field(
			&f.Inn,
			validation.Required,
			validation.Match(regexp.MustCompile(innPattern)),
		),
		validation.Field(&f.Name,
			validation.Required,
			validation.Match(regexp.MustCompile(namePattern))),
	)
}
