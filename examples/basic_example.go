package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/1saifj/reqpretty"
)

func main() {
	fmt.Println("🔬 Running Basic reqpretty Example...")

	// 1. Create a basic handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	// 2. Define minimal options
	opts := reqpretty.Options{
		IncludeRequest:      true,
		IncludeResponse:     true,
		IncludeResponseBody: true,
	}

	// 3. Wrap your handler with the middleware
	handler := reqpretty.DebugHandler(opts, nextHandler)

	// 4. Create a test request and recorder
	req := httptest.NewRequest(http.MethodGet, "http://example.com/basic", nil)
	rec := httptest.NewRecorder()

	// 5. Execute the request
	handler.ServeHTTP(rec, req)

	fmt.Println("\n✅ Basic Example Finished!")
}
