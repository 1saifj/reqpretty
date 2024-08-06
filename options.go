package reqpretty

// Options holds configuration for the debugger
type Options struct {
	IncludeRequest            bool
	IncludeRequestHeaders     bool
	IncludeRequestQueryParams bool
	IncludeRequestBody        bool
	IncludeResponse           bool
	IncludeResponseHeaders    bool
	IncludeResponseBody       bool
	SuccessEmoji              string
	ErrorEmoji                string
	ContextAttributes         []string
}
