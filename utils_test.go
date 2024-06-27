package reqpretty

import (
	"io"
	"strings"
	"testing"
)

func TestReadAndRestoreBody(t *testing.T) {
	bodyContent := "Hello, world!"
	body := io.NopCloser(strings.NewReader(bodyContent))
	result, err := readAndRestoreBody(body)
	if err != nil {
		t.Error("Expect no error when reading and restoring body")

	}

	if string(result) != bodyContent {
		t.Error("Expected body to be read and restored correctly")
	}

	restoredBody, _ := io.ReadAll(body)
	if string(restoredBody) != bodyContent {
		t.Error("Expected body to be restored correctly")
	}
}

func TestFormatBody(t *testing.T) {
	testCases := []struct {
		name     string
		body     []byte
		expected []string
	}{
		{
			name:     "JSON Body",
			body:     []byte(`{"message": "Hello, world!"}`),
			expected: []string{"  {", `    "message": "Hello, world!"`, "  }"},
		},
		{
			name:     "Non-JSON Body",
			body:     []byte("This is not JSON"),
			expected: []string{"This is not JSON"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatBody(tc.body)
			if !equalSlices(result, tc.expected) {
				t.Errorf("Expected:\n%v\nGot:\n%v", tc.expected, result)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
