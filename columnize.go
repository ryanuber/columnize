package columnize

import (
	"fmt"
	"strings"
)

// Returns a list of elements, each representing a single item which will
// belong to a column of output.
func getElementsFromLine(line string, delim string) []interface{} {
	elements := make([]interface{}, 0)
	for _, field := range strings.Split(line, delim) {
		elements = append(elements, strings.TrimSpace(field))
	}
	return elements
}

// Examines a list of strings and determines how wide each column should be
// considering all of the elements that need to be printed within it.
func getWidthsFromLines(lines []string, delim string) []int {
	var widths []int

	for _, line := range lines {
		elems := getElementsFromLine(line, delim)
		i := 0
		for _, elem := range elems {
			if len(widths) <= i {
				widths = append(widths, len(elem.(string)))
			} else if widths[i] < len(elem.(string)) {
				widths[i] = len(elem.(string))
			}
			i++
		}
	}
	return widths
}

// Given a set of column widths and the number of columns in the current line,
// returns a sprintf-style format string which can be used to print output
// aligned properly with other lines using the same widths set.
func getStringFormat(widths []int, columns int) string {
	var stringfmt string

	// Create the format string from the discovered widths
	for i := 0; i < columns && i < len(widths); i++ {
		if i == columns-1 {
			stringfmt += "%s\n"
		} else {
			stringfmt += fmt.Sprintf("%%-%ds", widths[i]+2)
		}
	}
	return stringfmt
}

// Columnize is the public-facing interface that takes a list of strings and a
// delimiter, and returns nicely aligned output.
func Columnize(input []string, delim string) string {
	var result string

	widths := getWidthsFromLines(input, delim)

	// Create the formatted output using the format string
	for _, line := range input {
		elems := getElementsFromLine(line, delim)
		stringfmt := getStringFormat(widths, len(elems))
		result += fmt.Sprintf(stringfmt, elems...)
	}
	return strings.TrimSpace(result)
}
