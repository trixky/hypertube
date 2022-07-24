package sanitizer

import (
	"strconv"
	"testing"
)

func TestSanitizeToken(t *testing.T) {
	tests := []struct {
		input          string
		error_expected bool
	}{
		// ------------------------- Failed expected
		{ // Input missing
			input:          "",
			error_expected: true,
		},
		{ // Input too short
			input:          "a",
			error_expected: true,
		},
		{ // One char missing
			input:          "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a",
			error_expected: true,
		},
		{ // One char too many
			input:          "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1a",
			error_expected: true,
		},
		{ // Corrupted char
			input:          "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a.",
			error_expected: true,
		},
		{ // Separators corrupted
			input:          "a1a1a1a1ba1a1ba1a1ba1a1ba1a1a1a1a1a.",
			error_expected: true,
		},
		// ------------------------- Success expected
		{
			input:          "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1",
			error_expected: false,
		},
		{
			input:          "a9c10773-3c34-47ec-91a1-b5d3a80ab38e",
			error_expected: false,
		},
		{
			input:          "df211c14-9980-4b8d-b8b4-ed612cc62a6d",
			error_expected: false,
		},
		{
			input:          "0f7825d6-00d2-4ca9-9aed-8c3f8bad71d5",
			error_expected: false,
		},
	}

	for _, test := range tests {
		if err := SanitizeToken(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeToken > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}
