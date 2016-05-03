Columnize
=========

Easy column-formatted output for golang

[![Build Status](https://travis-ci.org/ryanuber/columnize.svg)](https://travis-ci.org/ryanuber/columnize)

Columnize is a really small Go package that makes building CLI's a little bit
easier. In some CLI designs, you want to output a number similar items in a
human-readable way with nicely aligned columns. However, figuring out how wide
to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

```go
package main

import (
    "fmt"
    "github.com/CiscoCloud/columnize"
)

func main() {
    output := []string{
        "Name | Gender | Age",
        "Bob | Male | 38",
        "Sally | Female | 26",
    }
    result := columnize.SimpleFormat(output)
    fmt.Println(result)
}
```

As you can see, you just pass in a list of strings. And the result:

```
Name   Gender  Age
Bob    Male    38
Sally  Female  26
```

Columnize is tolerant of missing or empty fields, or even empty lines, so
passing in extra lines for spacing should show up as you would expect.

Configuration
=============

Columnize is configured using a `Config`, which can be obtained by calling the
`DefaultConfig()` method. You can then tweak the settings in the resulting
`Config`:

```
config := columnize.DefaultConfig()
config.Delim = "|"
config.Glue = "  "
config.Prefix = ""
config.Empty = ""
config.MaxWidth = []int{10, 0, 0}
config.OutputWidth = 80
```

* `Delim` is the string by which columns of **input** are delimited
* `Glue` is the string by which columns of **output** are delimited
* `Prefix` is a string by which each line of **output** is prefixed
* `Empty` is a string used to replace blank values found in output
* `MaxWidth` is an int slice specifying the maximum width of each column.
* `OutputWidth` is an int specifying the maximum width of an output line.

If MaxWidth or OutputWidth is specified and output exceeds the configured width, Columnize breaks a column at a word boundary and continues it on the next line.  See below for details.

You can then pass the `Config` in using the `Format` method (signature below) to
have text formatted to your liking.

Usage
=====

```go
SimpleFormat(intput []string) string

Format(input []string, config *Config) string
```

Controlling Output Width
========================
Output exceeding the width of the terminal window - particularly columnized output - can be difficult to read.  To address this, Columnize provides two configuration parameters for controlling output width.

* `MaxWidth` is an int slice specifying the maximum width of each column.  If the data for a column exceeds its maximum width, Columnize formats the column into two or more lines by breaking its data at a word boundary and continuing it onto the next line.  A zero or missing value for a MaxWidth element specifies that the corresponding column is uncontrolled (no maximum width).
* `OutputWidth` is an int value specifying the maximum width of the entire output line (including prefix and glue).  If data width exceeds this value, Columnize sets a MaxWidth for the rightmost uncontrolled column so that the output width satisfies the restriction.  You can specify `OutputWidth: columnize.AUTO` to use the actual width of the terminal window for OutputWidth.

For example:

    input := []string{
		"Column a | Column b | Column c",
		"xx | yy | zz",
		"some quite long data | some more data | even longer data for the last column",
		"this one will fit | a break | The quick brown fox jumps over the low lazy dog",
		"qq | rr | ss",
	}
	config := Config{MaxWidth: []int{10, 0, 15}}
	output := Format(input, &config)

results in the output:

    Column a    Column b        Column c
    xx          yy              zz
    some quite  some more data  even longer
    long data                   data for the
                                last column
    this one    a break         The quick brown
    will fit                    fox jumps over
                                the low lazy
                                dog
    qq          rr              ss

Specify OutputWidth to restrict the entire output line.  For example, the configuration:

    config := Config{
       OutputWidth: columnize.AUTO,
    }

causes the entire output line to fit in the terminal window.  Columnize modifies data lines exceeding the width of the window by setting the appropriate MaxWidth for the last column  of data.  Columnize adjusts the last uncontrolled column, so if you want it to adjust a column other than the last, specify an explicit MaxWidth for any columns to the right of the one you want Columnize to adjust.




