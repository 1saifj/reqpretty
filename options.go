package reqpretty

import (
	"log/slog"
	"strings"
)

// Options holds configuration for the debugger
type Options struct {
	IncludeRequest            bool
	IncludeRequestHeaders     bool
	IncludeRequestQueryParams bool
	IncludeRequestBody        bool
	IncludeResponse           bool
	IncludeResponseHeaders    bool
	IncludeResponseBody       bool
	Colorer                   Colorer
	EnableColor               bool // Add a boolean flag for coloring
	SuccessEmoji              string
	ErrorEmoji                string
}

// Colorer interface for custom coloring
type Colorer interface {
	Colorize(level slog.Level, msg string) string
}

// DefaultColorer provides basic colorization
type DefaultColorer struct{}

func (c *DefaultColorer) Colorize(level slog.Level, msg string) string {
	// Implement your custom coloring logic here based on the CuteColor example
	switch level {
	case slog.LevelInfo:
		if strings.Contains(msg, "REQUEST") || strings.Contains(msg, "RESPONSE") {
			return "\033[93m" + msg + "\033[0m" // Bright Yellow for titles
		}
		return "\033[95m" + msg + "\033[0m" // Bright Purple for messages
	case slog.LevelWarn:
		return "\033[33m" + msg + "\033[0m" // Yellow for Warning
	case slog.LevelError:
		return "\033[31m" + msg + "\033[0m" // Red for Error
	default:
		return msg // No color for other levels
	}
}

var config = Options{
	IncludeRequest:            true,
	IncludeRequestHeaders:     true,
	IncludeRequestQueryParams: true,
	IncludeRequestBody:        true,
	IncludeResponse:           true,
	IncludeResponseHeaders:    true,
	IncludeResponseBody:       true,
	Colorer:                   &DefaultColorer{},
	EnableColor:               true,
}
