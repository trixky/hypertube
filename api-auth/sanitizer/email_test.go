package sanitizer

import (
	"strconv"
	"testing"
)

func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		input          string
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing
			input:          "",
			error_expected: true,
		},
		{ // @ missing
			input:          "ab.c",
			error_expected: true,
		},
		{ // Domain missing #1
			input:          "a@b",
			error_expected: true,
		},
		{ // Domain missing #2
			input:          "a@b.",
			error_expected: true,
		},
		{ // @ and domain missing
			input:          "abc",
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input:          "a@b.c",
			error_expected: false,
		},
		{
			input:          "a.b.c@d.e",
			error_expected: false,
		},
		{
			input:          "ab@cd.ef",
			error_expected: false,
		},
		{
			input:          "chat@chat.chat",
			error_expected: false,
		},
		{
			input:          "chien@chien.chien",
			error_expected: false,
		},
		{
			input:          "chien.chien@chien.chien",
			error_expected: false,
		},
	}

	for _, test := range tests {
		if err := SanitizeEmail(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeEmail > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}
