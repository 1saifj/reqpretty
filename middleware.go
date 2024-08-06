package reqpretty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// DebugHandlerFunc is a function type for middleware
type DebugHandlerFunc func(opts Options, next http.Handler) http.Handler

// DebugHandler wraps an http.Handler with debug logging
func DebugHandler(opts Options, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Capture the request body and restore it for further processing
		reqBody, err := readAndRestoreBody(r.Body)
		if err != nil {
			slog.Error("Error reading request body", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Restore the body

		rec := newRecorder(w)
		next.ServeHTTP(rec, r)

		duration := time.Since(startTime)

		logRequest(r, reqBody, opts)
		logResponse(rec, duration, opts)
	})
}

// logRequest logs the request details
func logRequest(r *http.Request, reqBody []byte, opts Options) {
	if !opts.IncludeRequest {
		return
	}

	logger := slog.Default()
	logAttrs := []slog.Attr{
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
	}

	if opts.IncludeRequestHeaders {
		logAttrs = append(logAttrs, slog.Any("headers", r.Header))
	}
	if opts.IncludeRequestQueryParams && len(r.URL.Query()) > 0 {
		logAttrs = append(logAttrs, slog.Any("query_params", r.URL.Query()))
	}
	if opts.IncludeRequestBody && len(reqBody) > 0 {
		logAttrs = append(logAttrs, slog.Any("body", formatBody(reqBody)))
	}

	ctx := r.Context()
	logger = logger.With(convertAttrsToAny(extractContextAttributes(ctx, opts.ContextAttributes))...)

	logSection(logger, slog.LevelInfo, "⤴ REQUEST ⤴", logAttrs)
}

// logResponse logs the response details
func logResponse(rec *responseWriter, duration time.Duration, opts Options) {
	if !opts.IncludeResponse {
		return
	}

	logger := slog.Default()
	statusEmoji := opts.SuccessEmoji
	if rec.statusCode >= 400 {
		statusEmoji = opts.ErrorEmoji
	}

	logAttrs := []slog.Attr{
		slog.String("status", fmt.Sprintf("%d %s", rec.statusCode, http.StatusText(rec.statusCode))),
		slog.String("duration", duration.String()),
	}

	if opts.IncludeResponseHeaders {
		for name, values := range rec.Header() {
			logAttrs = append(logAttrs, slog.Any(name, values))
		}
	}
	if opts.IncludeResponseBody && len(rec.body) > 0 {
		logAttrs = append(logAttrs, slog.Any("body", formatBody(rec.body)))
	}

	logSection(logger, slog.LevelInfo, statusEmoji+" RESPONSE ⤵", logAttrs)
}

// logSection logs a section with a title and attributes
func logSection(logger *slog.Logger, level slog.Level, title string, attrs []slog.Attr) {
	logger.LogAttrs(nil, level, title, attrs...)
}

// readAndRestoreBody reads the request body and restores it for further processing
func readAndRestoreBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// formatBody formats the body for logging, handling JSON indentation
func formatBody(body []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
		return prettyJSON.String()
	}
	return string(body) // If not JSON, log as plain text
}

func convertAttrsToAny(attrs []slog.Attr) []any {
	converted := make([]any, len(attrs))
	for i, attr := range attrs {
		converted[i] = attr
	}
	return converted
}

func extractContextAttributes(ctx context.Context, attributes []string) []slog.Attr {
	var attrs []slog.Attr
	for _, attrName := range attributes {
		if attrValue := ctx.Value(attrName); attrValue != nil {
			attrs = append(attrs, slog.Any(attrName, attrValue))
		}
	}
	return attrs
}

// responseWriter wraps http.ResponseWriter to capture the status code and body
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func newRecorder(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.body = append(rw.body, p...) // Capture response body
	return rw.ResponseWriter.Write(p)
}
