package sanitizer

import (
	"errors"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SanitizeHexa64 sanitizes hexa64
func SanitizeHexa64(hexa64 string, label string) error {
	// https://stackoverflow.com/questions/336210/regular-expression-for-alphanumeric-and-underscores
	if len(hexa64) == 0 {
		return status.Errorf(codes.InvalidArgument, label+" missing")
	}
	if len(hexa64) != 64 {
		return status.Errorf(codes.InvalidArgument, label+" must be 64 characters long")
	}

	ok, err := regexp.MatchString("^[a-f0-9]*$", hexa64)

	if err != nil || !ok {
		return errors.New(label + " corrupted")
	}

	return nil
}
