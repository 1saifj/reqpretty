package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1saifj/reqpretty"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Configure reqpretty options
	opts := reqpretty.Options{
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

	// Configure the logger
	logger := &reqpretty.Logger{}
	reqpretty.Configure(logger)

	// Create routes
	mux := http.NewServeMux()

	// Simple GET endpoint
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":   "Hello, World!",
			"timestamp": time.Now().Unix(),
			"method":    r.Method,
		}
		json.NewEncoder(w).Encode(response)
	})

	// POST endpoint with JSON body
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Simulate creating user
		user.ID = 123
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	// Error endpoint for testing error logging
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	})

	// Wrap with reqpretty middleware
	loggedMux := reqpretty.DebugHandler(opts, mux)

	fmt.Println("🚀 Test server starting on :8080")
	fmt.Println("📝 Try these endpoints:")
	fmt.Println("   GET  http://localhost:8080/hello")
	fmt.Println("   POST http://localhost:8080/users")
	fmt.Println("   GET  http://localhost:8080/error")
	fmt.Println("👀 Watch the console for beautiful logs!")

	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
