package sanitizer

import (
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_token = "token"
)

func SanitizeToken(token string) error {
	// https://ihateregex.io/expr/uuid/
	if len(token) == 0 {
		return status.Errorf(codes.InvalidArgument, LABEL_token+" missing")
	}

	ok, err := regexp.MatchString("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", token)

	if err != nil {
		return status.Errorf(codes.InvalidArgument, LABEL_token+" corrupted (err):"+err.Error())
	}

	if !ok {
		return status.Errorf(codes.InvalidArgument, LABEL_token+" corrupted (ok):"+token)
	}

	if err != nil || !ok {
		return status.Errorf(codes.InvalidArgument, LABEL_token+" corrupted")
	}

	return nil
}
