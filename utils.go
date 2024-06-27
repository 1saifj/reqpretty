package reqpretty

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// readAndRestoreBody reads the request body and restores it for further processing.
func readAndRestoreBody(body io.ReadCloser) ([]byte, error) {
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	body = io.NopCloser(bytes.NewBuffer(buf))
	return buf, nil
}

// formatBody formats the body for logging, handling JSON indentation.
func formatBody(body []byte) []string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "  ", "  "); err == nil {
		return strings.Split(prettyJSON.String(), "\n")
	}
	return strings.Split(string(body), "\n") // If not JSON, log as plain text
}
