package sanitizer

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_username  = "username"
	LABEL_firstname = "firstname"
	LABEL_lastname  = "lastname"
)

func sanitizeName(name string, label string) error {
	if len(name) == 0 {
		return status.Errorf(codes.InvalidArgument, label+" missing")
	}
	if len(name) < 3 {
		return status.Errorf(codes.InvalidArgument, label+" too short")
	}
	if len(name) > 30 {
		return status.Errorf(codes.InvalidArgument, label+" too long")
	}

	return nil
}

func SanitizeUsername(username string) error {
	return sanitizeName(username, LABEL_username)
}

func SanitizeFirstname(username string) error {
	return sanitizeName(username, LABEL_firstname)
}

func SanitizeLastname(username string) error {
	return sanitizeName(username, LABEL_lastname)
}
