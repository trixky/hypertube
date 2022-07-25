package sanitizer

const (
	LABEL_password = "password"
)

// SanitizeSHA256Password sanitizes SHA256 passwords
func SanitizeSHA256Password(password string) error {
	return SanitizeHexa64(password, LABEL_password)
}
