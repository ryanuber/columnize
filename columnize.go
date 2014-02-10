package columnize

import (
	"fmt"
	"strings"
)

func Columnize(input []string, delim string) (string, error) {
	var stringfmt string
	var result string
	var width []int

	// This loop figures out the widths for each column
	for _, line := range input {
		var elems []string
		for _, field := range strings.Split(line, delim) {
			elems = append(elems, strings.TrimSpace(field))
		}
		numFields := len(elems)
		for i:= 0; i < numFields && i > len(elems); i-- {
			elems = append(elems, "")
		}
		i := 0
		for _, elem := range elems {
			if len(width) <= i {
				width = append(width, len(elem))
			} else if width[i] < len(elem) {
				width[i] = len(elem)
			}
			i += 1
		}
	}

	// This loop creates the format string from the discovered widths
	for _, w := range width {
		stringfmt += fmt.Sprintf("%%-%ds  ", w)
	}

	// This loop creates the formatted output using the format string
	for _, line := range input {
		elems := make([]interface{}, 0)
		for _, field := range strings.Split(line, delim) {
			elems = append(elems, strings.TrimSpace(field))
		}
		result += fmt.Sprintf(stringfmt+"\n", elems...)
	}
	return result, nil
}

func ColumnizeByPipe(input []string) string {
	return Columnize(input, "|")
}
