package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/pkg/protobuf/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
)

type SignUp struct {
	*core.SignUpMessage
}

func (f *SignUp) Validate() error {
	emailPattern := "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$"
	innPattern := fmt.Sprintf(`^\d{%d}$`, entities.InnLength)
	sliceErrorMessages := map[string]string{
		emailPattern: "Should contain email address",
		innPattern:   "Should contain 10 digits exactly",
	}
	return validation.ValidateStruct(f,
		validation.Field(
			&f.Login,
			validation.Required,
			validation.By(func(value interface{}) error {
				loginValue := value.(string)
				e := validation.Validate(loginValue, validation.Match(regexp.MustCompile(emailPattern)))
				if e != nil {
					return status.Error(codes.InvalidArgument, sliceErrorMessages[emailPattern])
				}
				return nil
			})),
		validation.Field(
			&f.Password,
			validation.Required,
			validation.Length(6, 120),
		),
		validation.Field(
			&f.Inn,
			validation.Required,
			validation.By(func(value interface{}) error {
				intValue := value.(uint64)
				stringValue := strconv.Itoa(int(intValue))
				e := validation.Validate(stringValue, validation.Match(regexp.MustCompile(innPattern)))
				if e != nil {
					return status.Error(codes.InvalidArgument, sliceErrorMessages[innPattern])
				}
				return nil
			}),
		),
		validation.Field(&f.Name,
			validation.Required,
			validation.Length(2, 120),
		),
		validation.Field(&f.Country,
			v.When(f.Country != "", is.CountryCode2),
		))
}
