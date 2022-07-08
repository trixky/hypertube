package sanitizer

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_password = "password"
)

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
	return SanitizeHexa64(password, LABEL_password)
}
