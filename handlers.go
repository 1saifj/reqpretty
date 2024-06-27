package reqpretty

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// DebugHandler wraps an http.Handler with debug logging
func DebugHandler(opts Options, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Capture the request body and restore it for further processing
		reqBody, err := readAndRestoreBody(r.Body)
		if err != nil {
			slog.Error("Error reading request body", "error", err)
		}

		rec := newRecorder(w)
		next.ServeHTTP(rec, r)

		duration := time.Since(startTime)

		logRequest(r, reqBody, opts)
		logResponse(rec, duration, opts)
	})
}

// logRequest logs the request details based on the given options
// handlers.go (updated logRequest function)

// logRequest logs the request details based on the given options
func logRequest(r *http.Request, reqBody []byte, opts Options) {
	logger := slog.Default()
	logSection(logger, opts.Colorer, slog.LevelInfo, LogSection{
		title:   "⤴ REQUEST ⤴",
		enabled: opts.IncludeRequest,
		content: func() []string {
			var lines []string
			lines = append(lines, fmt.Sprintf("%s %s", r.Method, r.URL))
			if opts.IncludeRequestHeaders {
				lines = append(lines, logHeaders(r.Header)...)
			}
			if opts.IncludeRequestQueryParams && len(r.URL.Query()) > 0 {
				lines = append(lines, "===== Query Parameters =====")
				for key, values := range r.URL.Query() {
					for _, value := range values {
						lines = append(lines, fmt.Sprintf("%s: %s", key, value))
					}
				}
			}
			if opts.IncludeRequestBody && len(reqBody) > 0 {
				lines = append(lines, formatBody(reqBody)...)
			}
			return lines
		},
	})
}

// Helper function to log headers
func logHeaders(headers http.Header) []string {
	var lines []string
	for name, values := range headers {
		lines = append(lines, fmt.Sprintf("%s: %s", name, strings.Join(values, ",")))
	}
	return lines
}

// logResponse logs the response details based on the given options
func logResponse(rec *responseWriter, duration time.Duration, opts Options) {
	logger := slog.Default()

	statusEmoji := opts.SuccessEmoji
	if rec.statusCode >= 400 {
		statusEmoji = opts.ErrorEmoji
	}

	logSection(logger, opts.Colorer, slog.LevelInfo, LogSection{
		title:   fmt.Sprintf("%s RESPONSE [%d/%s] [Time elapsed: %d ms]⤵", statusEmoji, rec.statusCode, http.StatusText(rec.statusCode), duration.Milliseconds()),
		enabled: opts.IncludeResponse,
		content: func() []string {
			lines := []string{}
			if opts.IncludeResponseHeaders {
				for name, values := range rec.Header() {
					lines = append(lines, fmt.Sprintf("%s: [%s]", name, strings.Join(values, ",")))
				}
			}
			if opts.IncludeResponseBody && len(rec.body) > 0 {
				lines = append(lines, formatBody(rec.body)...)
			}
			return lines
		},
	})
}
