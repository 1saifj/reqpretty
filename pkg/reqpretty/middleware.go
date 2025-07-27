// Package reqpretty provides HTTP middleware for beautiful request/response debugging
package reqpretty

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/1saifj/reqpretty/pkg/printer"
)

// DebugHandlerFunc is a function type for middleware
type DebugHandlerFunc func(opts Options, next http.Handler) http.Handler

// DebugHandler wraps an http.Handler with debug logging
func DebugHandler(opts Options, next http.Handler) http.Handler {
	if opts.Printer == nil {
		opts.Printer = printer.NewConsolePrinter()
	}
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

		defer func() {
			duration := time.Since(startTime)

			if rcv := recover(); rcv != nil {
				logPanic(rcv, opts.Printer)
				rec.WriteHeader(http.StatusInternalServerError)
			}

			// Always log the request and response
			logRequest(r, reqBody, opts)
			logResponse(rec, duration, opts)
		}()

		next.ServeHTTP(rec, r)
	})
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

// extractContextAttributes extracts specified attributes from context
func extractContextAttributes(ctx context.Context, attributes []string) []slog.Attr {
	var attrs []slog.Attr
	for _, attrName := range attributes {
		if attrValue := ctx.Value(attrName); attrValue != nil {
			attrs = append(attrs, slog.Any(attrName, attrValue))
		}
	}
	return attrs
}
