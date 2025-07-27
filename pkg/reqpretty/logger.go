package reqpretty

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/1saifj/reqpretty/pkg/printer"
)

// logRequest logs the request details in beautiful format
func logRequest(r *http.Request, reqBody []byte, opts Options) {
	if !opts.IncludeRequest {
		return
	}

	// Extract context attributes
	ctx := r.Context()
	contextAttrs := extractContextAttributes(ctx, opts.ContextAttributes)

	// Print request header
	printRequestHeader(r, opts.Printer)

	// Print context attributes if any
	if len(contextAttrs) > 0 {
		printContextAttributes(contextAttrs, opts.Printer)
	}

	// Print query parameters
	if opts.IncludeRequestQueryParams && len(r.URL.Query()) > 0 {
		opts.Printer.PrintTable(convertURLValuesToMap(r.URL.Query()), "Query Parameters")
	}

	// Print headers
	if opts.IncludeRequestHeaders {
		opts.Printer.PrintTable(convertHeadersToMap(r.Header), "Headers")
	}

	// Print body
	if opts.IncludeRequestBody && len(reqBody) > 0 {
		opts.Printer.PrintBody(reqBody, "Request Body")
	}
}

// logResponse logs the response details in beautiful format
func logResponse(rec *responseWriter, duration time.Duration, opts Options) {
	if !opts.IncludeResponse {
		return
	}

	// Print response header
	printResponseHeader(rec, duration, opts)

	// Print response headers
	if opts.IncludeResponseHeaders {
		opts.Printer.PrintTable(convertHeadersToMap(rec.Header()), "Response Headers")
	}

	// Print response body
	if opts.IncludeResponseBody && len(rec.body) > 0 {
		opts.Printer.PrintBody(rec.body, "Response Body")
	}

	// Add spacing after response
	fmt.Println()
}

// printRequestHeader prints a beautiful request header
func printRequestHeader(r *http.Request, printer printer.Printer) {
	method := r.Method
	url := r.URL.String()

	header := fmt.Sprintf("Request - %s", method)
	printer.PrintBox(header, url, "blue")
}

// printResponseHeader prints a beautiful response header
func printResponseHeader(rec *responseWriter, duration time.Duration, opts Options) {
	statusEmoji := opts.SuccessEmoji
	statusColor := "green"

	if rec.statusCode >= 400 {
		statusEmoji = opts.ErrorEmoji
		statusColor = "red"
	}

	if statusEmoji == "" {
		if rec.statusCode >= 400 {
			statusEmoji = "❌"
		} else {
			statusEmoji = "✅"
		}
	}

	status := fmt.Sprintf("%d %s", rec.statusCode, http.StatusText(rec.statusCode))
	timeStr := duration.String()

	header := fmt.Sprintf("%s Response - Status: %s - Time: %s", statusEmoji, status, timeStr)
	opts.Printer.PrintBox(header, "", statusColor)
}

// printContextAttributes prints context attributes in a table
func printContextAttributes(attrs []slog.Attr, printer printer.Printer) {
	if len(attrs) == 0 {
		return
	}

	contextMap := make(map[string]interface{})
	for _, attr := range attrs {
		contextMap[attr.Key] = attr.Value.Any()
	}
	printer.PrintTable(contextMap, "Context Attributes")
}

// logPanic prints panic details in a beautiful box
func logPanic(rcv interface{}, printer printer.Printer) {
	errorMsg := fmt.Sprintf("💥 PANIC RECOVERED 💥\n\nError: %v\n", rcv)
	// Just a portion of the stack to avoid huge logs
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	if len(lines) > 15 {
		lines = lines[:15]
		lines = append(lines, "...")
	}
	stackMsg := fmt.Sprintf("Stack Trace (truncated):\n%s", strings.Join(lines, "\n"))
	content := errorMsg + "\n" + stackMsg

	printer.PrintBox("PANIC", content, "red")
}

// Helper functions to convert types
func convertHeadersToMap(headers http.Header) map[string]interface{} {
	result := make(map[string]interface{})
	for key, values := range headers {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}
	return result
}

func convertURLValuesToMap(values map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, vals := range values {
		if len(vals) == 1 {
			result[key] = vals[0]
		} else {
			result[key] = vals
		}
	}
	return result
}
