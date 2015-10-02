package columnize

import (
	"fmt"
	"strings"
	"unicode"
)

type Config struct {
	// The string by which the lines of input will be split.
	Delim string

	// The string by which columns of output will be separated.
	Glue string

	// The string by which columns of output will be prefixed.
	Prefix string

	// A replacement string to replace empty fields
	Empty string

	// Maximum width of each field
	MaxWidth []int

	// Index of a negative value within config.MaxWidth
	negIndex int
}

// Returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Delim:    "|",
		Glue:     "  ",
		Prefix:   "",
		MaxWidth: []int{},
	}
}

// Returns a list of elements, each representing a single item which will
// belong to a column of output.
func getElementsFromLine(config *Config, line string) []interface{} {
	elements := make([]interface{}, 0)
	for _, field := range strings.Split(line, config.Delim) {
		value := strings.TrimSpace(field)
		if value == "" && config.Empty != "" {
			value = config.Empty
		}
		elements = append(elements, value)
	}
	return elements
}

// Examines a list of strings and determines how wide each column should be
// considering all of the elements that need to be printed within it.
func getWidthsFromLines(config *Config, lines []string) []int {
	var widths []int

	for _, line := range lines {
		elems := getElementsFromLine(config, line)
		for i := 0; i < len(elems); i++ {
			l := len(elems[i].(string))
			if i < len(config.MaxWidth) && config.MaxWidth[i] > 0 && config.MaxWidth[i] < l {
				l = config.MaxWidth[i]
			}
			if len(widths) <= i {
				widths = append(widths, l)
			} else if widths[i] < l {
				widths[i] = l
			}
		}
	}

	// If one of the columns has a negative width specification, set its width
	// so that the entire output line has a width of the absolute value of the spec

	if config.negIndex >= 0 && config.negIndex < len(widths) {
		maxOutputWidth := 0 - config.MaxWidth[config.negIndex]
		if config.negIndex == 0 && len(config.MaxWidth) == 1 && len(widths) > 1 {
			config.negIndex = len(widths) - 1 // Apply single negative value to last field
		}
		totalLineWidth := len(config.Prefix) + len(config.Glue)*(len(widths)-1)
		for i, width := range widths {
			if i != config.negIndex {
				totalLineWidth += width
			}
		}
		if maxOutputWidth > totalLineWidth {
			widths[config.negIndex] = maxOutputWidth - totalLineWidth
		} else {
			config.negIndex = -1
		}
	}
	return widths
}

// Given a set of column widths and the number of columns in the current line,
// returns a sprintf-style format string which can be used to print output
// aligned properly with other lines using the same widths set.
func (c *Config) getStringFormat(widths []int, columns int) string {
	// Start with the prefix, if any was given.
	stringfmt := c.Prefix

	// Create the format string from the discovered widths
	for i := 0; i < columns && i < len(widths); i++ {
		if i == columns-1 {
			stringfmt += "%s\n"
		} else {
			stringfmt += fmt.Sprintf("%%-%ds%s", widths[i], c.Glue)
		}
	}
	return stringfmt
}

// MergeConfig merges two config objects together and returns the resulting
// configuration. Values from the right take precedence over the left side.
func MergeConfig(a, b *Config) *Config {
	var result Config = *a
	result.negIndex = -1

	// Return quickly if either side was nil
	if a == nil || b == nil {
		return &result
	}

	if b.Delim != "" {
		result.Delim = b.Delim
	}
	if b.Glue != "" {
		result.Glue = b.Glue
	}
	if b.Prefix != "" {
		result.Prefix = b.Prefix
	}
	if b.Empty != "" {
		result.Empty = b.Empty
	}
	if len(b.MaxWidth) > 0 {
		for i, maxWidth := range b.MaxWidth {
			if maxWidth < 0 {
				// Negative width adjusts the column width so that the width of
				// the entire output line is the abs value of the value specified.
				// Only the first such specification is significant; others are ignored.
				if result.negIndex < 0 {
					result.negIndex = i
				} else {
					maxWidth = 0
				}
			}
			if i < len(result.MaxWidth) {
				result.MaxWidth[i] = maxWidth
			} else {
				result.MaxWidth = append(result.MaxWidth, maxWidth)
			}
		}
	}

	return &result
}

// Format is the public-facing interface that takes either a plain string
// or a list of strings and returns nicely aligned output.
func Format(lines []string, config *Config) string {
	var result string

	conf := MergeConfig(DefaultConfig(), config)
	widths := getWidthsFromLines(conf, lines)

	// Create the formatted output using the format string
	for _, line := range lines {
		elems := getElementsFromLine(conf, line)
		extensionLineElems := []string{}
		isStillDataToFormat := true
		for isStillDataToFormat {
			isStillDataToFormat = truncateToWidth(&elems, &extensionLineElems, widths)
			stringfmt := conf.getStringFormat(widths, len(elems))
			result += fmt.Sprintf(stringfmt, elems...)
		}
	}

	// Remove trailing newline without removing leading/trailing space
	if n := len(result); n > 0 && result[n-1] == '\n' {
		result = result[:n-1]
	}

	return result
}

// Convenience function for using Columnize as easy as possible.
func SimpleFormat(lines []string) string {
	return Format(lines, nil)
}

// Truncate any elements exceeding their maximum width, and save their remaining
// data for an extension line.
func truncateToWidth(elems *[]interface{}, extensionLineElems *[]string, widths []int) (isStillDataToFormat bool) {

	// If this an extension line, make its list of elements current

	if len(*extensionLineElems) > 0 {
		for i, elem := range *extensionLineElems {
			(*elems)[i] = elem
		}
		*extensionLineElems = []string{}
	}

	// Examine each element to determine if it exceeds its maximum allowed width.
	// If so, truncate it at the closest whitespace to the limit and save its remaining
	// data for the next extension line.

	for i, elem := range *elems {
		stringElem := strings.TrimSpace(fmt.Sprintf("%s", elem))
		if len(stringElem) > widths[i] {
			isStillDataToFormat = true
			splitPoint := widths[i]
			for ; splitPoint > 0; splitPoint-- {
				if unicode.IsSpace(rune(stringElem[splitPoint])) {
					break
				}
			}
			if splitPoint == 0 {
				splitPoint = widths[i]
			}
			(*elems)[i] = strings.TrimSpace(stringElem[:splitPoint])
			if len(*extensionLineElems) == 0 {
				(*extensionLineElems) = make([]string, len(*elems))
			}
			(*extensionLineElems)[i] = strings.TrimSpace(stringElem[splitPoint:])
		}

	}
	return
}
