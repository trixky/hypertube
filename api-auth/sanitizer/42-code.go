package sanitizer

const (
	LABEL_42_code = "code"
)

// Sanitize42Code sanitizes 42 codes
func Sanitize42Code(code string) error {
	return SanitizeHexa64(code, LABEL_42_code)
}
