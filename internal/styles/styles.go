// Package styles provides styling utilities for the CLI tool
package styles

import (
	"strings"
)

// ANSI color codes
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"

	// Colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Box drawing characters
	TopLeft     = "┌"
	TopRight    = "┐"
	BottomLeft  = "└"
	BottomRight = "┘"
	Horizontal  = "─"
	Vertical    = "│"
)

type Style struct {
	color   string
	bgColor string
	bold    bool
	padding int
	margin  int
	border  bool
	width   int
}

func NewStyle() Style {
	return Style{}
}

func (s Style) Foreground(color string) Style {
	s.color = color
	return s
}

func (s Style) Bold(bold bool) Style {
	s.bold = bold
	return s
}

func (s Style) Border(border bool) Style {
	s.border = border
	return s
}

func (s Style) Padding(padding int) Style {
	s.padding = padding
	return s
}

func (s Style) Render(text string) string {
	result := text

	// Apply text styling
	if s.bold {
		result = Bold + result
	}
	if s.color != "" {
		result = s.color + result
	}

	// Add reset at the end
	result += Reset

	// Apply padding
	if s.padding > 0 {
		padding := strings.Repeat(" ", s.padding)
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			lines[i] = padding + line + padding
		}
		result = strings.Join(lines, "\n")
	}

	// Apply border
	if s.border {
		result = s.addBorder(result)
	}

	return result
}

func (s Style) addBorder(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	// Calculate the width of the box
	maxWidth := 0
	for _, line := range lines {
		cleanLine := removeAnsiCodes(line)
		if len(cleanLine) > maxWidth {
			maxWidth = len(cleanLine)
		}
	}

	// Create bordered content
	var result []string

	// Top border
	topBorder := TopLeft + strings.Repeat(Horizontal, maxWidth) + TopRight
	result = append(result, topBorder)

	// Content lines
	for _, line := range lines {
		cleanLine := removeAnsiCodes(line)
		padding := maxWidth - len(cleanLine)
		paddedLine := Vertical + line + strings.Repeat(" ", padding) + Vertical
		result = append(result, paddedLine)
	}

	// Bottom border
	bottomBorder := BottomLeft + strings.Repeat(Horizontal, maxWidth) + BottomRight
	result = append(result, bottomBorder)

	return strings.Join(result, "\n")
}

func removeAnsiCodes(text string) string {
	result := ""
	inEscape := false
	for i, char := range text {
		if char == '\033' && i+1 < len(text) && text[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape && (char == 'm' || char == 'K' || char == 'J') {
			inEscape = false
			continue
		}
		if !inEscape {
			result += string(char)
		}
	}
	return result
}

// Helper functions for common styles
func ColoredText(text, color string) string {
	return color + text + Reset
}

func BoldText(text string) string {
	return Bold + text + Reset
}
