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
		SuccessEmoji:              "✅",
		ErrorEmoji:                "❌",
		ContextAttributes:         []string{"request_id", "user_id"},
	}

	t.Run("test request with 200 status code and body", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, world!"))
		})
		handler := DebugHandler(opts, nextHandler)

		reqBody := []byte(`{"key":"value"}`)
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo?bar=baz", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "request_id", "12345saif"))
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

	t.Run("test request with 400 status code and body", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		})
		handler := DebugHandler(opts, nextHandler)

		reqBody := []byte(`{"key":"value"}`)
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo?bar=baz", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Long-Header", "This is a very long header value that should wrap to multiple lines to test the wrapping functionality in the logging box")
		req = req.WithContext(context.WithValue(req.Context(), "request_id", "12345saif"))
		req = req.WithContext(context.WithValue(req.Context(), "user_id", "user-1"))

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}
		expectedBody := "Bad Request"
		if string(body) != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", string(body), expectedBody)
		}
	})

	t.Run("test panic recovery", func(t *testing.T) {
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something went terribly wrong")
		})
		handler := DebugHandler(opts, panicHandler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/panic", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("test request with no body", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("could not read request body: %v", err)
			}
			if len(body) != 0 {
				t.Errorf("expected empty request body but got %d bytes", len(body))
			}
			w.WriteHeader(http.StatusOK)
		})
		handler := DebugHandler(opts, nextHandler)

		req := httptest.NewRequest(http.MethodPost, "http://example.com/no-body", nil) // No body provided
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test with non-json body", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		handler := DebugHandler(opts, nextHandler)

		reqBody := []byte("this is plain text")
		req := httptest.NewRequest(http.MethodPost, "http://example.com/text", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "text/plain") // Set content type
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test with multi-value header", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		handler := DebugHandler(opts, nextHandler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/multi-header", nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept", "application/xml")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("test with 204 no content response", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent) // No body should be written for 204
		})
		handler := DebugHandler(opts, nextHandler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/no-content", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		if status := resp.StatusCode; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}
		if len(body) != 0 {
			t.Errorf("expected empty response body but got %d bytes", len(body))
		}
	})
}
