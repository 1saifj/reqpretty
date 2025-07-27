package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ANSI color codes and styling
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"

	// Colors
	Blue    = "\033[34m"
	Green   = "\033[32m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors
	BrightBlue   = "\033[94m"
	BrightGreen  = "\033[92m"
	BrightRed    = "\033[91m"
	BrightYellow = "\033[93m"
	BrightCyan   = "\033[96m"

	// Box drawing characters
	TopLeft     = "┌"
	TopRight    = "┐"
	BottomLeft  = "└"
	BottomRight = "┘"
	Horizontal  = "─"
	Vertical    = "│"
	TeeDown     = "┬"
	TeeUp       = "┴"
)

// ConsolePrinter implements Printer interface for console output
type ConsolePrinter struct{}

// NewConsolePrinter creates a new console printer
func NewConsolePrinter() *ConsolePrinter {
	return &ConsolePrinter{}
}

// PrintBox prints text in a beautiful bordered box
func (p *ConsolePrinter) PrintBox(header, text, color string) {
	colorCode := p.getColorCode(color)
	content := header
	if text != "" {
		content += "\n" + text
	}

	lines := strings.Split(content, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Add padding
	maxWidth += 4

	// Top border
	fmt.Printf("%s%s%s%s%s\n", colorCode, TopLeft, strings.Repeat(Horizontal, maxWidth), TopRight, Reset)

	// Content lines
	for _, line := range lines {
		padding := maxWidth - len(line) - 2
		fmt.Printf("%s%s %s%s %s%s\n", colorCode, Vertical, Reset+line, strings.Repeat(" ", padding), colorCode, Vertical+Reset)
	}

	// Bottom border
	fmt.Printf("%s%s%s%s%s\n", colorCode, BottomLeft, strings.Repeat(Horizontal, maxWidth), BottomRight, Reset)
	fmt.Println()
}

// PrintTable prints a map as a beautiful table
func (p *ConsolePrinter) PrintTable(data map[string]interface{}, header string) {
	if len(data) == 0 {
		return
	}

	// Print table header
	fmt.Printf("%s%s%s %s %s%s\n", BrightCyan, Bold, header, Reset, BrightCyan, Reset)

	// Calculate max key width
	maxKeyWidth := 0
	for key := range data {
		if len(key) > maxKeyWidth {
			maxKeyWidth = len(key)
		}
	}
	maxKeyWidth += 2 // Add padding

	// Calculate max value width
	maxValueWidth := 0
	for _, value := range data {
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) > maxValueWidth {
			maxValueWidth = len(valueStr)
		}
	}
	maxValueWidth += 2 // Add padding

	// Top border
	fmt.Printf("%s%s%s%s%s%s\n",
		BrightCyan, TopLeft,
		strings.Repeat(Horizontal, maxKeyWidth),
		TeeDown,
		strings.Repeat(Horizontal, maxValueWidth),
		TopRight+Reset)

	// Data rows
	for key, value := range data {
		valueStr := fmt.Sprintf("%v", value)
		keyPadding := maxKeyWidth - len(key) - 1
		valuePadding := maxValueWidth - len(valueStr) - 1

		fmt.Printf("%s%s%s %s%s%s%s %s%s%s%s\n",
			BrightCyan, Vertical, Reset,
			key, strings.Repeat(" ", keyPadding),
			BrightCyan, Vertical, Reset,
			valueStr, strings.Repeat(" ", valuePadding),
			BrightCyan, Vertical)
	}

	// Bottom border
	fmt.Printf("%s%s%s%s%s%s\n",
		BrightCyan, BottomLeft,
		strings.Repeat(Horizontal, maxKeyWidth),
		TeeUp,
		strings.Repeat(Horizontal, maxValueWidth),
		BottomRight+Reset)
	fmt.Println()
}

// PrintBody prints formatted body content
func (p *ConsolePrinter) PrintBody(body []byte, header string) {
	formattedBody := p.formatBodyPretty(body)

	// Print header
	fmt.Printf("%s%s%s %s %s%s\n", BrightYellow, Bold, header, Reset, BrightYellow, Reset)

	lines := strings.Split(formattedBody, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 4 // Add padding

	// Top border
	fmt.Printf("%s%s%s%s%s\n", BrightYellow, TopLeft, strings.Repeat(Horizontal, maxWidth), TopRight, Reset)

	// Content lines
	for _, line := range lines {
		padding := maxWidth - len(line) - 2
		fmt.Printf("%s%s %s%s %s%s\n", BrightYellow, Vertical, Reset+line, strings.Repeat(" ", padding), BrightYellow, Vertical+Reset)
	}

	// Bottom border
	fmt.Printf("%s%s%s%s%s\n", BrightYellow, BottomLeft, strings.Repeat(Horizontal, maxWidth), BottomRight, Reset)
	fmt.Println()
}

// getColorCode returns the ANSI color code for a color name
func (p *ConsolePrinter) getColorCode(color string) string {
	switch strings.ToLower(color) {
	case "blue":
		return BrightBlue
	case "green":
		return BrightGreen
	case "red":
		return BrightRed
	case "yellow":
		return BrightYellow
	case "cyan":
		return BrightCyan
	default:
		return White
	}
}

// formatBodyPretty formats the body for pretty printing
func (p *ConsolePrinter) formatBodyPretty(body []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "    "); err == nil {
		return prettyJSON.String()
	}
	return string(body) // If not JSON, return as plain text
}
