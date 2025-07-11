package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/1saifj/reqpretty"
)

func main() {
	fmt.Println("🔬 Running Advanced reqpretty Example...")
	// Create a new router
	mux := http.NewServeMux()
	// This handler will be wrapped by our middleware
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		// Simulate some work
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "user created successfully"}`))
	})
	// This handler will demonstrate an error response
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "error", "message": "invalid input provided"}`))
	})
	// Define all available options for the middleware
	opts := reqpretty.Options{
		IncludeRequest:            true,
		IncludeRequestQueryParams: true,
		IncludeRequestHeaders:     true,
		IncludeRequestBody:        true,
		IncludeResponse:           true,
		IncludeResponseHeaders:    true,
		IncludeResponseBody:       true,
		ContextAttributes:         []string{"request_id", "session_id", "user_id"},
		SuccessEmoji:              "🎉",
		ErrorEmoji:                "🔥",
	}
	// Wrap the router with the reqpretty middleware
	handler := reqpretty.DebugHandler(opts, mux)
	// --- Simulate a successful API call ---
	fmt.Println("\n\n🚀 Simulating a SUCCESSFUL request to /user...")
	// Create a test request with context, headers, and body
	reqSuccessBody := []byte(`{"username": "testuser", "email": "test@example.com"}`)
	reqSuccess := httptest.NewRequest(http.MethodPost, "http://localhost:8080/user?source=test", bytes.NewReader(reqSuccessBody))
	reqSuccess.Header.Set("Content-Type", "application/json")
	reqSuccess.Header.Set("X-Custom-Header", "my-custom-value")
	// Add context values that match the `ContextAttributes` in options
	ctxSuccess := context.WithValue(reqSuccess.Context(), "request_id", "req-12345")
	ctxSuccess = context.WithValue(ctxSuccess, "session_id", "sess-abcde")
	ctxSuccess = context.WithValue(ctxSuccess, "user_id", "usr-zyxw")
	reqSuccess = reqSuccess.WithContext(ctxSuccess)
	// Create a recorder to capture the response
	recSuccess := httptest.NewRecorder()
	// Execute the request
	handler.ServeHTTP(recSuccess, reqSuccess)
	// --- Simulate a failed API call ---
	fmt.Println("\n\n💥 Simulating a FAILED request to /error...")
	// Create another test request
	reqErrorBody := []byte(`{"username": "invalid-user"}`)
	reqError := httptest.NewRequest(http.MethodPost, "http://localhost:8080/error", bytes.NewReader(reqErrorBody))
	reqError.Header.Set("Content-Type", "application/json")
	// Add different context for this request
	ctxError := context.WithValue(reqError.Context(), "request_id", "req-67890")
	reqError = reqError.WithContext(ctxError)
	// Create a new recorder and execute the request
	recError := httptest.NewRecorder()
	handler.ServeHTTP(recError, reqError)
	fmt.Println("\n\n✅ Example Finished!")
}
