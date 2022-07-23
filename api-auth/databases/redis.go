package databases

import "strings"

const (
	REDIS_EXTERNAL_google            = "google"
	REDIS_SEPARATOR                  = "."
	REDIS_PATTERN_KEY_token          = "token"
	REDIS_PATTERN_KEY_password_token = "password_token"
)

// ErrorIsDuplication checks if an sql error is an duplication error
func ErrorIsDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}
