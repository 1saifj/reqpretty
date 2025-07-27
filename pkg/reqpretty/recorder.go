package reqpretty

import "net/http"

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
