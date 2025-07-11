package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1saifj/reqpretty"
)

func main() {
	fmt.Println("🔬 Advanced reqpretty testing...")

	// Test 1: Minimal logging (production-like)
	fmt.Println("\n1️⃣ Testing minimal logging...")
	testMinimalLogging()

	// Test 2: Headers only
	fmt.Println("\n2️⃣ Testing headers-only logging...")
	testHeadersOnly()

	// Test 3: Context attributes
	fmt.Println("\n3️⃣ Testing context attributes...")
	testContextAttributes()

	fmt.Println("\n✅ All advanced tests completed!")
}

func testMinimalLogging() {
	opts := reqpretty.Options{
		IncludeRequest:      true,
		IncludeResponse:     true,
		IncludeRequestBody:  false, // Skip body for privacy
		IncludeResponseBody: false,
		SuccessEmoji:        "✓",
		ErrorEmoji:          "✗",
	}

	handler := reqpretty.DebugHandler(opts, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Minimal response"))
	}))

	// Simulate request
	req, _ := http.NewRequest("GET", "http://example.com/minimal", nil)
	rec := &mockResponseWriter{}
	handler.ServeHTTP(rec, req)
}

func testHeadersOnly() {
	opts := reqpretty.Options{
		IncludeRequest:         true,
		IncludeRequestHeaders:  true,
		IncludeResponse:        true,
		IncludeResponseHeaders: true,
		SuccessEmoji:           "🎯",
		ErrorEmoji:             "💥",
	}

	handler := reqpretty.DebugHandler(opts, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Headers test"))
	}))

	// Simulate request with headers
	req, _ := http.NewRequest("POST", "http://example.com/headers", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("User-Agent", "reqpretty-advanced-test")
	rec := &mockResponseWriter{}
	handler.ServeHTTP(rec, req)
}

func testContextAttributes() {
	opts := reqpretty.Options{
		IncludeRequest:    true,
		IncludeResponse:   true,
		SuccessEmoji:      "🚀",
		ContextAttributes: []string{"trace_id", "user_id", "session_id"},
	}

	handler := reqpretty.DebugHandler(opts, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Context test successful",
		})
	}))

	// Simulate request with context
	req, _ := http.NewRequest("GET", "http://example.com/context", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "trace_id", "trace-789")
	ctx = context.WithValue(ctx, "user_id", "user-456")
	ctx = context.WithValue(ctx, "session_id", "sess-123")
	req = req.WithContext(ctx)

	rec := &mockResponseWriter{}
	handler.ServeHTTP(rec, req)
}

// Mock response writer for testing
type mockResponseWriter struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.body = append(m.body, data...)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(code int) {
	m.statusCode = code
}
