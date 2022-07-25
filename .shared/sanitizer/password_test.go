package sanitizer

import (
	"strconv"
	"testing"
)

func TestSanitizeSHA256Password(t *testing.T) {
	for _, test := range hex64_tests {
		if err := SanitizeSHA256Password(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeSHA256Password > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}
