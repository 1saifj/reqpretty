package reqpretty

import (
	"fmt"
	"log/slog"
	"strings"
)

type LogSection struct {
	title   string
	enabled bool
	content func() []string
}

// logSection logs a section of the request/response details
// logSection logs a section of the request/response details
func logSection(logger *slog.Logger, colorer Colorer, level slog.Level, section LogSection) {
	if !config.EnableColor { // If coloring is disabled, set colorer to nil
		colorer = nil
	}
	logger.Info(colorer.Colorize(level, drawTitle(section.title)))
	for _, line := range section.content() {
		logger.Info(colorer.Colorize(level, drawMessage(line)))
	}
}

// local draw title with box
func drawTitle(title string) (box string) {
	// box elements
	topright := "╮"
	topleft := "╭"
	w := "─"
	h := "│"
	bottomright := "╯"
	bottomleft := "╰"

	// calculate the length of line
	length := len(title)
	// draw line
	line := strings.Repeat(w, length+2)

	// print title box
	box = fmt.Sprintf("%v%v%v\n", topleft, line, topright)
	box += fmt.Sprintf("%v %v %v\n", h, title, h)
	box += fmt.Sprintf("%v%v%v", bottomleft, line, bottomright)
	return
}

// local draw message
func drawMessage(message string) (msg string) {
	msg = fmt.Sprintf("🭬 %v", message)
	return
}
