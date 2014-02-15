package columnize

import (
	"strings"
	"testing"
)

func TestColumnizedOutput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}
	output := Columnize(input, "|")
	expected := strings.TrimSpace(strings.Join([]string{
		"Column A  Column B  Column C",
		"x         y         z",
	}, "\n"))
	if expected != output {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyLinesOutput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"",
		"x | y | z",
	}
	output := Columnize(input, "|")
	expected := strings.TrimSpace(strings.Join([]string{
		"Column A  Column B  Column C",
		"",
		"x         y         z",
	}, "\n"))
	if expected != output {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}
