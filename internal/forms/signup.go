package forms

import "go-user-microservice/pkg/protobuf/user"

type SignUp struct{
	*user.SignUpMessage
}
