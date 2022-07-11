package sanitizer

import (
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LABEL_id = "email"
)

func SanitizeID(id string) (int64, error) {
	if len(id) == 0 {
		return 0, status.Errorf(codes.InvalidArgument, LABEL_id+" missing")
	}

	id_integer, err := strconv.Atoi(id)

	if err != nil || id_integer < 0 {
		return 0, status.Errorf(codes.InvalidArgument, LABEL_id+" corrupted")
	}

	return int64(id_integer), nil
}
