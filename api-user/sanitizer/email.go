package sanitizer

import (
	"net/mail"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_email = "email"
)

func SanitizeEmail(email string) error {
	if len(email) == 0 {
		return status.Errorf(codes.InvalidArgument, LABEL_email+" missing")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return status.Errorf(codes.InvalidArgument, LABEL_email+" corrupted")
	}

	return nil
}
