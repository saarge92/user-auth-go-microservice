package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DatabaseError(err error) error {
	return err
}

func CustomDatabaseError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}
