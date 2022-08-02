package tools

import (
	"regexp"
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

var (
	ellipsis            = "..."
	minWidthForEllipsis = len(ellipsis) + 2
	lineRE 				= regexp.MustCompile(`(?m)^`)
	ws 					= regexp.MustCompile(`\s+`)
)

func Indent(s, indent string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}

	return lineRE.ReplaceAllLiteralString(s, indent)
}

func ReplaceExcessiveWhitespace(s string) string {
	return ws.ReplaceAllString(strings.TrimSpace(s), " ")
}

// DisplayWidth calculates what the rendered width of a string may be
func DisplayWidth(s string) int {
	return ansi.PrintableRuneWidth(s)
}

// Truncate shortens a string to fit the maximum display width
func Truncate(maxWidth int, s string) string {
	w := DisplayWidth(s)
	if w <= maxWidth {
		return s
	}

	tail := ""
	if maxWidth >= minWidthForEllipsis {
		tail = ellipsis
	}

	r := truncate.StringWithTail(s, uint(maxWidth), tail)
	if DisplayWidth(r) < maxWidth {
		r += " "
	}

	return r
}

// TruncateColumn replaces the first new line character with an ellipsis
// and shortens a string to fit the maximum display width
func TruncateColumn(maxWidth int, s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i] + ellipsis
	}

	return Truncate(maxWidth, s)
}
