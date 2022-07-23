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

// sanitizeName sanitizes names
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

// SanitizeUsername sanitizes username
func SanitizeUsername(username string) error {
	return sanitizeName(username, LABEL_username)
}

// SanitizeFirstname sanitizes firstname
func SanitizeFirstname(firstname string) error {
	return sanitizeName(firstname, LABEL_firstname)
}

// SanitizeLastname sanitizes lastname
func SanitizeLastname(lastname string) error {
	return sanitizeName(lastname, LABEL_lastname)
}
