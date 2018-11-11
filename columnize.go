package columnize

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Config can be used to tune certain parameters which affect the way
// in which Columnize will format output text.
type Config struct {
	// The string by which the lines of input will be split.
	Delim string

	// The string by which columns of output will be separated.
	Glue string

	// The string by which columns of output will be prefixed.
	Prefix string

	// A replacement string to replace empty fields
	Empty string
}

// DefaultConfig returns a *Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Delim:  "|",
		Glue:   "  ",
		Prefix: "",
		Empty:  "",
	}
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

	return &result
}

// stringFormat, given a set of column widths and the number of columns in
// the current line, returns a sprintf-style format string which can be used
// to print output aligned properly with other lines using the same widths set.
func stringFormat(c *Config, widths []int, columns int) string {
	// Create the buffer with an estimate of the length
	buf := bytes.NewBuffer(make([]byte, 0, (6+len(c.Glue))*columns))

	// Start with the prefix, if any was given. The buffer will not return an
	// error so it does not need to be handled
	buf.WriteString(c.Prefix)

	// Create the format string from the discovered widths
	for i := 0; i < columns && i < len(widths); i++ {
		if i == columns-1 {
			buf.WriteString("%s\n")
		} else {
			fmt.Fprintf(buf, "%%-%ds%s", widths[i], c.Glue)
		}
	}
	return buf.String()
}

// elementsFromLine returns a list of elements, each representing a single
// item which will belong to a column of output.
func elementsFromLine(config *Config, line string) []interface{} {
	seperated := strings.Split(line, config.Delim)
	elements := make([]interface{}, len(seperated))
	for i, field := range seperated {
		value := strings.TrimSpace(field)

		// Apply the empty value, if configured.
		if value == "" && config.Empty != "" {
			value = config.Empty
		}
		elements[i] = value
	}
	return elements
}

// runeLen calculates the number of visible "characters" in a string
func runeLen(s string) int {
	l := 0
	insideConsoleColor := false
	for _, c := range s {
		if c == '\x1b' { // start of esc color sequence
			insideConsoleColor = true
			continue
		}
		if insideConsoleColor {
			if c == 'm' { // end of esc color sequence
				insideConsoleColor = false
			}
			continue
		}

		l++
		if unicode.Is(unicode.Scripts["Han"], rune(c)) {
			l++
		}

	}

	return l
}

// widthsFromLines examines a list of strings and determines how wide each
// column should be considering all of the elements that need to be printed
// within it.
func widthsFromLines(config *Config, lines []string) []int {
	widths := make([]int, 0, 8)

	for _, line := range lines {
		elems := elementsFromLine(config, line)
		for i := 0; i < len(elems); i++ {
			l := runeLen(elems[i].(string))
			if len(widths) <= i {
				widths = append(widths, l)
			} else if widths[i] < l {
				widths[i] = l
			}
		}
	}

	if os.Getenv("DEBUG_COL") == "1" {
		fmt.Println(widths)
	}
	return widths
}

// Format is the public-facing interface that takes a list of strings and
// returns nicely aligned column-formatted text.
func Format(lines []string, config *Config) string {
	conf := MergeConfig(DefaultConfig(), config)
	widths := widthsFromLines(conf, lines)

	// Estimate the buffer size
	glueSize := len(conf.Glue)
	var size int
	for _, w := range widths {
		size += w + glueSize
	}
	size *= len(lines)

	// Create the buffer
	buf := bytes.NewBuffer(make([]byte, 0, size))

	// Create a cache for the string formats
	fmtCache := make(map[int]string, 16)

	// Create the formatted output using the format string
	for _, line := range lines {
		elems := elementsFromLine(conf, line)

		// Get the string format using cache
		numElems := len(elems)
		stringfmt, ok := fmtCache[numElems]
		if !ok {
			stringfmt = stringFormat(conf, widths, numElems)
			fmtCache[numElems] = stringfmt
		}

		if false {
			// 原作者的逻辑
			// fmt.Fprintf("%-21s  %s", "我们", "12.21.212.1"), 如果中文会有问题
			fmt.Fprintf(buf, stringfmt, elems...)
		} else {
			// 手动解决
			fmt.Fprintf(buf, conf.Prefix)
			for col, elem := range elems {
				s := elem.(string)
				for _, c := range s {
					fmt.Fprintf(buf, string(c))
				}

				if col == len(elems)-1 {
					// 最后一列
					fmt.Fprintf(buf, "\n")
				} else {
					// space padding
					fmt.Fprintf(buf, strings.Repeat(" ", widths[col]-runeLen(s)))
					// glue padding
					fmt.Fprintf(buf, conf.Glue)
				}

			}
		}
	}

	// Get the string result
	result := buf.String()

	// Remove trailing newline without removing leading/trailing space
	if n := len(result); n > 0 && result[n-1] == '\n' {
		result = result[:n-1]
	}

	return result
}

// SimpleFormat is a convenience function to format text with the defaults.
func SimpleFormat(lines []string) string {
	return Format(lines, nil)
}
