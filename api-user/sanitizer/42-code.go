package sanitizer

const (
	LABEL_42_code = "code"
)

func Sanitize42Code(code string) error {
	return SanitizeHexa64(code, LABEL_42_code)
}
