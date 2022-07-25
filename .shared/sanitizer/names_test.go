package sanitizer

import (
	"strconv"
	"testing"
)

var name_tests = []struct {
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
		input:          "aa",
		error_expected: true,
	},
	// ------------------------- Success expected
	{
		input:          "ncolomer",
		error_expected: false,
	},
	{
		input:          "mabois",
		error_expected: false,
	},
	{
		input:          "oggy",
		error_expected: false,
	},
}

func TestSanitizeName(t *testing.T) {
	for _, test := range name_tests {
		if err := sanitizeName(test.input, "nothing"); (err != nil) != test.error_expected {
			t.Fatalf("sanitizeName > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}

func TestSanitizeUsername(t *testing.T) {
	for _, test := range name_tests {
		if err := SanitizeUsername(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeUsername > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}

func TestSanitizeFirstname(t *testing.T) {
	for _, test := range name_tests {
		if err := SanitizeFirstname(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeFirstname > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}

func TestSanitizeLastname(t *testing.T) {
	for _, test := range name_tests {
		if err := SanitizeLastname(test.input); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeLastname > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}
