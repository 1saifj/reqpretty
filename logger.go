package reqpretty

import (
	"log/slog"
	"path/filepath"
)

// Default is the default *Logger for the default handler.
var Default = &Logger{}

// Configure configures the logger by the given options.
func Configure(logger *Logger) {
	if logger == nil {
		logger = Default
	}
	customHandler := NewHandler(logger)
	slog.SetDefault(slog.New(customHandler))
}

type Logger struct {
	slog.Handler
}

func (o *Logger) clone() *Logger {
	return &Logger{
		Handler: o.Handler.WithAttrs(nil), // Make sure to clone the attributes
	}
}

// defaultReplaceAttrFunc is a function that can be used to replace the value of an attribute.
// remove time key and source prefix
var defaultReplaceAttrFunc = func(groups []string, a slog.Attr) slog.Attr {
	// Remove time attributes from the log.
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	// Remove source filepath prefix attributes from the log.
	if a.Key == slog.SourceKey {
		source, ok := a.Value.Any().(*slog.Source)
		if !ok {
			return slog.Attr{}
		}
		if source != nil {
			source.File = filepath.Base(source.File)
		}
	}
	return a
}
