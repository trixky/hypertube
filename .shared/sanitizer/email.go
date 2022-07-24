package sanitizer

import (
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_email = "email"
)

func SanitizeEmail(email string) error {
	// https://regexr.com/3e48o // (modified)

	if len(email) == 0 {
		return status.Errorf(codes.InvalidArgument, LABEL_email+" missing")
	}

	ok, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{1,63}$`, email)

	if err != nil || !ok {
		return status.Errorf(codes.InvalidArgument, LABEL_email+" corrupted")
	}

	return nil
}
