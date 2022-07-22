package utils

import (
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// ------------------------- Success expected
		{
			input:    "22280755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			expected: "9ba9862e6c992abe652e575881931b22a29676b328d6cc00d04f1659888d40a1",
		},
		{
			input:    "cd180755e9747a2b40ec92502dbb76f612049fb0f7a2926216e2bdcfa849f368",
			expected: "c47b37adeba3497e1de2005bf3d0f8c28806a813527ce9f3cb917ae896b41cd4",
		},
	}

	for _, test := range tests {
		result := EncryptPassword(test.input)
		if test.expected != result {
			t.Fatalf("EncryptPassword > expected: " + test.expected + " | result: " + result)
		}
	}
}
