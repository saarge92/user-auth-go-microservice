package forms

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"go-user-microservice/internal/entites"
	"go-user-microservice/pkg/protobuf/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
)

type SignUp struct {
	*user.SignUpMessage
	userExistRule validation.RuleFunc
	userInnRule   validation.RuleFunc
}

func NewSignUpForm(
	request *user.SignUpMessage,
	userExistRule validation.RuleFunc,
	userInnRule validation.RuleFunc,
) *SignUp {
	return &SignUp{
		request,
		userExistRule,
		userInnRule,
	}
}

func (f *SignUp) Validate() error {
	emailPattern := "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$"
	innPattern := fmt.Sprintf(`^\d{%d}$`, entites.InnLength)
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
			}),
			validation.By(f.userExistRule)),
		validation.Field(
			&f.Password,
			validation.Required,
			validation.Length(6, 120),
		),
		validation.Field(
			&f.Inn,
			validation.Required,
			validation.By(func(value interface{}) error {
				intValue := value.(uint32)
				stringValue := strconv.Itoa(int(intValue))
				e := validation.Validate(stringValue, validation.Match(regexp.MustCompile(innPattern)))
				if e != nil {
					return status.Error(codes.InvalidArgument, sliceErrorMessages[innPattern])
				}
				return nil
			}),
			validation.By(f.userInnRule),
		),
		validation.Field(&f.Name,
			validation.Required,
			validation.Length(2, 120),
		),
	)
}
