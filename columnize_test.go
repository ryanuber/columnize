package columnize

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestListOfStringsInput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyLinesOutput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestLeadingSpacePreserved(t *testing.T) {
	input := []string{
		"| Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "   Column B  Column C\n"
	expected += "x  y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestColumnWidthCalculator(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Longer than A | Longer than B | Longer than C",
		"short | short | short",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A       Column B       Column C\n"
	expected += "Longer than A  Longer than B  Longer than C\n"
	expected += "short          short          short"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestVariedInputSpacing(t *testing.T) {
	input := []string{
		"Column A       |Column B|    Column C",
		"x|y|          z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestUnmatchedColumnCounts(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Value A | Value B",
		"Value A | Value B | Value C | Value D",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "Value A   Value B\n"
	expected += "Value A   Value B   Value C   Value D"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateDelimiter(t *testing.T) {
	input := []string{
		"Column | A % Column | B % Column | C",
		"Value A % Value B % Value C",
	}

	config := DefaultConfig()
	config.Delim = "%"
	output := Format(input, config)

	expected := "Column | A  Column | B  Column | C\n"
	expected += "Value A     Value B     Value C"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateSpacingString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Glue = "    "
	output := Format(input, config)

	expected := "Column A    Column B    Column C\n"
	expected += "x           y           z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestSimpleFormat(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	output := SimpleFormat(input)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternatePrefixString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Prefix = "  "
	output := Format(input, config)

	expected := "  Column A  Column B  Column C\n"
	expected += "  x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyFieldReplacement(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | | z",
	}

	config := DefaultConfig()
	config.Empty = "<none>"
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         <none>    z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyConfigValues(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := Config{}
	output := Format(input, &config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestMergeConfig(t *testing.T) {
	conf1 := &Config{Delim: "a", Glue: "a", Prefix: "a", Empty: "a"}
	conf2 := &Config{Delim: "b", Glue: "b", Prefix: "b", Empty: "b"}
	conf3 := &Config{Delim: "c", Prefix: "c"}

	m := MergeConfig(conf1, conf2)
	if m.Delim != "b" || m.Glue != "b" || m.Prefix != "b" || m.Empty != "b" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, conf3)
	if m.Delim != "c" || m.Glue != "a" || m.Prefix != "c" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, nil)
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, &Config{})
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}
}

func TestMaxWidth(t *testing.T) {
	input := []string{
		"Column a | Column b | Column c",
		"xx | yy | zz",
		"some quite long data | some more data | even longer data for the last column",
		"this one will fit | a break | The quick brown fox jumps over the low lazy dog",
		"qq | rr | ss",
	}
	config := Config{MaxWidth: []int{10, 0, 15}}
	output := Format(input, &config)
	expected := "Column a    Column b        Column c\n"
	expected += "xx          yy              zz\n"
	expected += "some quite  some more data  even longer\n"
	expected += "long data                   data for the\n"
	expected += "                            last column\n"
	expected += "this one    a break         The quick brown\n"
	expected += "will fit                    fox jumps over\n"
	expected += "                            the low lazy\n"
	expected += "                            dog\n"
	expected += "qq          rr              ss"

	if output != expected {
		for i, c := range output {
			expectedChar := " "
			if i < len(expected) {
				expectedChar = string(expected[i])
			}
			if c != rune(expectedChar[0]) {
				nearStart := int(math.Max(0, float64(i-4)))
				nearEnd := int(math.Min(float64(i+4), float64(len(output))))
				near := strings.Replace(output[nearStart:nearEnd], "\n", " ", -1)
				fmt.Printf("TestMaxWidth difference at column %d near \"%s\": got(%s) expected(%s)\n", i, near, string(c), string(expectedChar))
			}
		}
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestNegativeMaxWidth(t *testing.T) {
	input := []string{
		"Column a | Column b | Column c",
		"xx | yy | zz",
		"some quite long data | some more data | even longer data for the last column",
		"this one will fit | a break | The quick brown fox jumps over the low lazy dog",
		"qq | rr | ss",
	}
	config := Config{MaxWidth: []int{0, 0, -60}}
	output := Format(input, &config)
	expected := "Column a              Column b        Column c\n"
	expected += "xx                    yy              zz\n"
	expected += "some quite long data  some more data  even longer data for\n"
	expected += "                                      the last column\n"
	expected += "this one will fit     a break         The quick brown fox\n"
	expected += "                                      jumps over the low\n"
	expected += "                                      lazy dog\n"
	expected += "qq                    rr              ss"

	if output != expected {
		for i, c := range output {
			expectedChar := " "
			if i < len(expected) {
				expectedChar = string(expected[i])
			}
			if c != rune(expectedChar[0]) {
				nearStart := int(math.Max(0, float64(i-4)))
				nearEnd := int(math.Min(float64(i+4), float64(len(output))))
				near := strings.Replace(output[nearStart:nearEnd], "\n", " ", -1)
				fmt.Printf("TestNegativeMaxWidth difference at column %d near \"%s\": got(%s) expected(%s)\n", i, near, string(c), string(expectedChar))
			}
		}
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}
