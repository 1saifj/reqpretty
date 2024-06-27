package reqpretty

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDebugHandler(t *testing.T) {
	opts := Options{
		IncludeRequest:            true, // or adjust as needed
		IncludeRequestHeaders:     true,
		IncludeRequestQueryParams: true,
		IncludeRequestBody:        true,
		IncludeResponse:           true,
		IncludeResponseHeaders:    true,
		IncludeResponseBody:       true,
		SuccessEmoji:              "✅",
		ErrorEmoji:                "❌",
		Colorer:                   &DefaultColorer{}, // Use the default colorer
		EnableColor:               true,              // or false to disable colors
	}

	handler := DebugHandler(opts, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"message": "Hello, world!"})
		if err != nil {
			return
		}
	}))

	reqBody := bytes.NewBufferString(`{"name": "Alice"}`)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var logOutput bytes.Buffer
	handlerOptions := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(&logOutput, &handlerOptions))
	originalLogger := slog.Default()
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	handler.ServeHTTP(w, req)

	output := logOutput.String()

	expectedLogElements := []string{
		"POST /",
		"Content-Type: application/json",
		"name: Alice",
		"RESPONSE [200/OK]",
		"message: Hello, world!",
	}

	for _, element := range expectedLogElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected log to contain '%s', but it didn't. Log output:\n%s", element, output)
		}
	}
}
