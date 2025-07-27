// Package printer provides interfaces and implementations for formatting output
package printer

// Printer defines the interface for formatting and printing debug output
type Printer interface {
	// PrintBox prints text in a bordered box with specified color
	PrintBox(header, content, color string)

	// PrintTable prints a map as a formatted table
	PrintTable(data map[string]interface{}, header string)

	// PrintBody prints formatted body content
	PrintBody(body []byte, header string)
}
