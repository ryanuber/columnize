package columnize

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
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

	// Maximum width of output; set to AUTO to use actual console width
	OutputWidth int

	// Maximum width of each column
	MaxWidth []int
}

const (
	AUTO = -1
)

var (
	defaultConfig Config = Config{
		Delim:       "|",
		Glue:        "  ",
		Prefix:      "",
		OutputWidth: -999,
		MaxWidth:    []int{},
	}
)

// Returns a Config with default values.
func DefaultConfig() *Config {
	var defaultConfigCopy Config = defaultConfig
	return &defaultConfigCopy
}

// Sets the default Config
func SetDefaultConfig(config Config) {
	defaultConfig = config
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
			lenElem := len(elems[i].(string))
			if i < len(config.MaxWidth) {
				if config.MaxWidth[i] < 0 {
					fmt.Printf("Columnize: negative MaxWidth value not supported - please use OutputWidth\n")
				} else if config.MaxWidth[i] > 0 && config.MaxWidth[i] < lenElem {
					lenElem = config.MaxWidth[i]
				}
			}
			if len(widths) <= i {
				widths = append(widths, lenElem)
			} else if widths[i] < lenElem {
				widths[i] = lenElem
			}
		}
	}

	// Get output width restriction

	outputWidth := config.OutputWidth
	if outputWidth == AUTO {
		var e error
		if outputWidth, e = GetConsoleWidth(); e != nil {
			fmt.Printf("Unable to set AUTO OutputWidth: %s\n", e.Error())
		}
	}

	// If the output width is restricted and the output line will exceed that width,
	// attempt to meet the restriction by adjusting the width of the rightmost
	// unrestricted column, or the rightmost column if all columns are restricted.

	if outputWidth > 0 {
		totalLineWidth := len(config.Prefix) + len(config.Glue)*(len(widths)-1)
		for _, width := range widths {
			totalLineWidth += width
		}
		if totalLineWidth > outputWidth {
			adjIndex := -1
			for i := len(widths) - 1; i >= 0; i-- {
				if i >= len(config.MaxWidth) || config.MaxWidth[i] <= 0 {
					adjIndex = i
					break
				}
			}
			if adjIndex < 0 {
				adjIndex = len(widths) - 1
			}
			adjustedWidth := outputWidth - (totalLineWidth - widths[adjIndex])
			if adjustedWidth > 0 {
				widths[adjIndex] = adjustedWidth
			}
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
	if b.OutputWidth >= 0 || b.OutputWidth == AUTO {
		result.OutputWidth = b.OutputWidth
	}
	if len(b.MaxWidth) > 0 {
		result.MaxWidth = b.MaxWidth
	}

	return &result
}

// Format is the public-facing interface that takes either a plain string
// or a list of strings and returns nicely aligned output.
func Format(lines []string, config *Config) string {
	var result string

	conf := MergeConfig(&defaultConfig, config)
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

// Get the width of the console
func GetConsoleWidth() (width int, e error) {
	var rxGetWidth *regexp.Regexp
	var command []string
	if runtime.GOOS == "windows" {
		command = []string{"mode", "con"}
		rxGetWidth = regexp.MustCompile("Columns:\\s*(\\d+)")
	} else {
		command = []string{"stty", "size"}
		rxGetWidth = regexp.MustCompile("\\d+\\s+(\\d+)")
	}
	cmd := exec.Command(command[0], command[1])
	cmd.Stdin = os.Stdin
	if out, err := cmd.Output(); err != nil {
		fmt.Printf("Unable to get console width: %s (%s)\n", err.Error(), out)
	} else if match := rxGetWidth.FindSubmatch(out); match != nil {
		width, e = strconv.Atoi(string(match[1]))
	} else {
		e = fmt.Errorf("not in %s %s output\n%s", command[0], command[1], out)
	}
	return
}
