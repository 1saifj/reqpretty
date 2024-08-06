package reqpretty

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BufferLogger struct {
	buffer bytes.Buffer
}

func (l *BufferLogger) Info(msg string, attrs ...slog.Attr) {
	l.buffer.WriteString(msg + "\n")
	for _, attr := range attrs {
		l.buffer.WriteString(attr.Key + ": " + attr.Value.String() + "\n")
	}
}

func (l *BufferLogger) String() string {
	return l.buffer.String()
}

func TestDebugHandler(t *testing.T) {
	opts := Options{
		IncludeRequest:            true,
		IncludeRequestHeaders:     true,
		IncludeRequestQueryParams: true,
		IncludeRequestBody:        true,
		IncludeResponse:           true,
		IncludeResponseHeaders:    true,
		IncludeResponseBody:       true,
		SuccessEmoji:              "✔️",
		ErrorEmoji:                "❗",
		ContextAttributes:         []string{"request_id", "user_id"},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	handler := DebugHandler(opts, nextHandler)

	t.Run("test request and response logging", func(t *testing.T) {
		reqBody := []byte(`{"key":"value"}`)
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo?bar=baz", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "request_id", "12345"))
		req = req.WithContext(context.WithValue(req.Context(), "user_id", "user-1"))

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}
		expectedBody := "Hello, world!"
		if string(body) != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", string(body), expectedBody)
		}
	})
}
