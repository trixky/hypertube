package sanitizer

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_email     = "email"
	LABEL_username  = "username"
	LABEL_firstname = "firstname"
	LABEL_lastname  = "lastname"
	LABEL_password  = "password"
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

func SanitizeName(name string, label string) error {
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

func SanitizePassword(password string) error {
	if len(password) == 0 {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" missing")
	}
	if len(password) < 8 {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" too short")
	}
	if len(password) > 30 {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" too long")
	}

	const numeric = "0123456789" // [0-9]

	if !strings.ContainsAny(password, numeric) {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" must contain at least one number")
	}

	const alpha = "qwertyuioplkjhgfdsazxcvbnm" // [a-z]

	if !strings.ContainsAny(password, alpha) {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" must contain at least one lowercase character")
	}
	if !strings.ContainsAny(password, strings.ToUpper(alpha)) { // [A-Z]
		return status.Errorf(codes.InvalidArgument, LABEL_password+" must contain at least one uppercase character")
	}

	const special = " !@#$%^&*()-=_+[]{}\\|'\";:/?.>,<`~"

	if !strings.ContainsAny(password, special) { // [!@#$...]
		return status.Errorf(codes.InvalidArgument, LABEL_password+" must contain at least one special character")
	}

	return nil
}

func SanitizeSHA256Password(password string) error {
	// https://stackoverflow.com/questions/336210/regular-expression-for-alphanumeric-and-underscores
	if len(password) == 0 {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" missing")
	}
	if len(password) != 64 {
		return status.Errorf(codes.InvalidArgument, LABEL_password+" must be 64 characters long")
	}

	ok, err := regexp.MatchString("^[a-zA-Z0-9_]*$", password)

	fmt.Println("apres regex:")
	fmt.Println(ok)
	fmt.Println(err)

	if err != nil {
		return err
	}
	if !ok {
		return errors.New("contains corrupted characters")
	}

	return nil
}
