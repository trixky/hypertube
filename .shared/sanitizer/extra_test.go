package sanitizer

import (
	"strconv"
	"testing"
)

var hex64_tests = []struct {
	input          string
	error_expected bool
}{
	// ------------------------- Failed expected
	{ //Input missing
		input:          "",
		error_expected: true,
	},
	{ //Input too short
		input:          "a",
		error_expected: true,
	},
	{ // One char missing
		input:          "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a",
		error_expected: true,
	},
	{ // One char too many
		input:          "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a",
		error_expected: true,
	},
	{ // Corrupted char
		input:          "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a.",
		error_expected: true,
	},
	// ------------------------- Success expected
	{
		input:          "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
		error_expected: false,
	},
	{
		input:          "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
		error_expected: false,
	},
	{
		input:          "cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
		error_expected: false,
	},
}

func TestSanitizeHexa64(t *testing.T) {
	for _, test := range hex64_tests {
		if err := SanitizeHexa64(test.input, "nothing"); (err != nil) != test.error_expected {
			t.Fatalf("SanitizeHexa64 > expected an error: " + strconv.FormatBool(test.error_expected))
		}
	}
}
