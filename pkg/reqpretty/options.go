package reqpretty

import "github.com/1saifj/reqpretty/pkg/printer"

// Options configures the debug middleware behavior
type Options struct {
	// Request logging options
	IncludeRequest            bool
	IncludeRequestHeaders     bool
	IncludeRequestBody        bool
	IncludeRequestQueryParams bool

	// Response logging options
	IncludeResponse        bool
	IncludeResponseHeaders bool
	IncludeResponseBody    bool

	// Context attributes to log
	ContextAttributes []string

	// Custom emojis for status indication
	SuccessEmoji string
	ErrorEmoji   string

	// Printer interface for customizable output formatting
	Printer printer.Printer
}

// DefaultOptions returns sensible default options
func DefaultOptions() Options {
	return Options{
		IncludeRequest:            true,
		IncludeRequestHeaders:     true,
		IncludeRequestBody:        true,
		IncludeRequestQueryParams: true,
		IncludeResponse:           true,
		IncludeResponseHeaders:    true,
		IncludeResponseBody:       true,
		ContextAttributes:         []string{},
		SuccessEmoji:              "✅",
		ErrorEmoji:                "❌",
		Printer:                   printer.NewConsolePrinter(),
	}
}
