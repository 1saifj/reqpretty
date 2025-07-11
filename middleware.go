package reqpretty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
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

		defer func() {
			duration := time.Since(startTime)

			if rcv := recover(); rcv != nil {
				logPanic(rcv)
				rec.WriteHeader(http.StatusInternalServerError) // Write header to the actual response
			}

			// Always log the request and response
			logRequest(r, reqBody, opts)
			logResponse(rec, duration, opts)
		}()

		next.ServeHTTP(rec, r)
	})
}

// Constants for beautiful closed box logging
const (
	maxWidth    = 90
	tabStep     = "    "
	topLeft     = "╔"
	topRight    = "╗"
	bottomLeft  = "╚"
	bottomRight = "╝"
	horizontal  = "═"
	vertical    = "║"
)

// calculateVisualWidth calculates the visible width of a string, accounting for known double-width emojis.
// This is a workaround to avoid adding external dependencies for terminal-aware string width calculation.
// It may not work for all emojis or in all terminals.
func calculateVisualWidth(s string) int {
	width := 0
	for _, r := range s {
		switch r {
		case '✅', '❌', '🎯', '🚀', '🔬':
			width += 2 // These emojis often render as double-width
		default:
			width += 1
		}
	}
	return width
}

// logRequest logs the request details in beautiful format
func logRequest(r *http.Request, reqBody []byte, opts Options) {
	if !opts.IncludeRequest {
		return
	}

	// Extract context attributes
	ctx := r.Context()
	contextAttrs := extractContextAttributes(ctx, opts.ContextAttributes)

	// Print request header
	printRequestHeader(r)

	// Print context attributes if any
	if len(contextAttrs) > 0 {
		printContextAttributes(contextAttrs)
	}

	// Print query parameters
	if opts.IncludeRequestQueryParams && len(r.URL.Query()) > 0 {
		printMapAsTable(convertURLValuesToMap(r.URL.Query()), "Query Parameters")
	}

	// Print headers
	if opts.IncludeRequestHeaders {
		printMapAsTable(convertHeadersToMap(r.Header), "Headers")
	}

	// Print body
	if opts.IncludeRequestBody && len(reqBody) > 0 {
		printBody(reqBody, "Request Body")
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
		printMapAsTable(convertHeadersToMap(rec.Header()), "Response Headers")
	}

	// Print response body
	if opts.IncludeResponseBody && len(rec.body) > 0 {
		printBody(rec.body, "Response Body")
	}

	// Add spacing after response
	fmt.Println()
}

// printRequestHeader prints a beautiful request header
func printRequestHeader(r *http.Request) {
	method := r.Method
	url := r.URL.String()

	header := fmt.Sprintf("Request - %s", method)
	printBoxed(header, url)
}

// printResponseHeader prints a beautiful response header
func printResponseHeader(rec *responseWriter, duration time.Duration, opts Options) {
	statusEmoji := opts.SuccessEmoji
	if rec.statusCode >= 400 {
		statusEmoji = opts.ErrorEmoji
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
	printBoxed(header, "")
}

// printContextAttributes prints context attributes in a table
func printContextAttributes(attrs []slog.Attr) {
	if len(attrs) == 0 {
		return
	}

	contextMap := make(map[string]interface{})
	for _, attr := range attrs {
		contextMap[attr.Key] = attr.Value.Any()
	}
	printMapAsTable(contextMap, "Context Attributes")
}

// printBoxed prints text in a beautiful closed box
func printBoxed(header, text string) {
	fmt.Println()
	printTopLine()
	printContentLine(header)
	if text != "" {
		printContentLine(text)
	}
	printBottomLine()
}

// printMapAsTable prints a map as a beautiful closed table
func printMapAsTable(data map[string]interface{}, header string) {
	if len(data) == 0 {
		return
	}

	printTopLine()
	printContentLine(header)
	for key, value := range data {
		printKVClosed(key, value)
	}
	printBottomLine()
}

// printTopLine prints the top border of a closed box
func printTopLine() {
	line := strings.Repeat(horizontal, 86) // Fixed width for consistent alignment
	fmt.Printf("%s%s%s\n", topLeft, line, topRight)
}

// printBottomLine prints the bottom border of a closed box
func printBottomLine() {
	line := strings.Repeat(horizontal, 86) // Fixed width for consistent alignment
	fmt.Printf("%s%s%s\n", bottomLeft, line, bottomRight)
}

// printContentLine prints a line with content inside closed borders
func printContentLine(content string) {
	maxContentWidth := 84 // 90 - 6 (borders and spaces)
	contentLength := calculateVisualWidth(content)
	if contentLength <= maxContentWidth {
		padding := maxContentWidth - contentLength
		fmt.Printf("%s %s%s %s\n", vertical, content, strings.Repeat(" ", padding), vertical)
	} else {
		// Content too long, wrap it
		printWrappedContent(content)
	}
}

// printWrappedContent prints long content with proper wrapping in closed borders
func printWrappedContent(content string) {
	maxContentWidth := 84 // Fixed width for consistent alignment
	currentLine := ""
	currentWidth := 0

	for _, r := range []rune(content) {
		runeWidth := 1
		switch r {
		case '✅', '❌', '🎯', '🚀', '🔬':
			runeWidth = 2
		}

		if currentWidth+runeWidth > maxContentWidth {
			padding := maxContentWidth - currentWidth
			fmt.Printf("%s %s%s %s\n", vertical, currentLine, strings.Repeat(" ", padding), vertical)
			currentLine = string(r)
			currentWidth = runeWidth
		} else {
			currentLine += string(r)
			currentWidth += runeWidth
		}
	}

	if currentWidth > 0 {
		padding := maxContentWidth - currentWidth
		fmt.Printf("%s %s%s %s\n", vertical, currentLine, strings.Repeat(" ", padding), vertical)
	}
}

// printBody prints request/response body in a beautiful closed format
func printBody(body []byte, header string) {
	printTopLine()
	printContentLine(header)
	printEmptyLine()

	formatted := formatBodyPretty(body)
	lines := strings.Split(formatted, "\n")
	for _, line := range lines {
		printBodyLine(line)
	}

	printEmptyLine()
	printBottomLine()
}

// printKVClosed prints a key-value pair in a closed box format
func printKVClosed(key string, value interface{}) {
	valueStr := fmt.Sprintf("%v", value)
	content := fmt.Sprintf("%s: %s", key, valueStr)
	printContentLine(content)
}

// printEmptyLine prints an empty line inside closed borders
func printEmptyLine() {
	padding := strings.Repeat(" ", 84) // Fixed width for consistent alignment
	fmt.Printf("%s %s %s\n", vertical, padding, vertical)
}

// printBodyLine prints a body content line inside closed borders
func printBodyLine(line string) {
	content := tabStep + line
	printContentLine(content)
}

// logPanic prints panic details in a beautiful box.
func logPanic(rcv interface{}) {
	printTopLine()
	printContentLine("💥 PANIC RECOVERED 💥")
	printEmptyLine()
	printWrappedContent(fmt.Sprintf("Error: %v", rcv))
	printEmptyLine()
	// Just a portion of the stack to avoid huge logs
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	if len(lines) > 15 {
		lines = lines[:15]
		lines = append(lines, "...")
	}
	printWrappedContent(fmt.Sprintf("Stack Trace (truncated):\n%s", strings.Join(lines, "\n")))
	printBottomLine()
}

// printKV prints a key-value pair (legacy function)
func printKV(key string, value interface{}) {
	printKVClosed(key, value)
}

// printLine prints a horizontal line (legacy function - kept for compatibility)
func printLine(start string) {
	line := strings.Repeat(horizontal, maxWidth)
	fmt.Printf("%s%s%s\n", start, line, bottomRight)
}

// printBlock prints a text block with proper line wrapping (legacy function)
func printBlock(text string) {
	printWrappedContent(text)
}

// formatBodyPretty formats the body for pretty printing
func formatBodyPretty(body []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", tabStep); err == nil {
		return prettyJSON.String()
	}
	return string(body) // If not JSON, return as plain text
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

// formatBody formats the body for logging, handling JSON indentation (legacy function)
func formatBody(body []byte) string {
	return formatBodyPretty(body)
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
