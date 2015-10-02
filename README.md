Columnize
=========

Easy column-formatted output for golang

[![Build Status](https://travis-ci.org/ryanuber/columnize.svg)](https://travis-ci.org/ryanuber/columnize)

Columnize is a really small Go package that makes building CLI's a little bit
easier. In some CLI designs, you want to output a number similar items in a
human-readable way with nicely aligned columns. However, figuring out how wide
to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

    package main

    import (
        "fmt"
        "github.com/tooda02/columnize"
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

As you can see, you just pass in a list of strings. And the result:

    Name   Gender  Age
    Bob    Male    38
    Sally  Female  26

Columnize is tolerant of missing or empty fields, or even empty lines, so
passing in extra lines for spacing should show up as you would expect.

Configuration
=============

Columnize is configured using a `Config`, which can be obtained by calling the
`DefaultConfig()` method. You can then tweak the settings in the resulting
`Config`:

    config := columnize.DefaultConfig()
    config.Delim = "|"
    config.Glue = "  "
    config.Prefix = ""
    config.Empty = ""
    config.MaxWidth = []int{10, 0, -80}

* `Delim` is the string by which columns of **input** are delimited
* `Glue` is the string by which columns of **output** are delimited
* `Prefix` is a string by which each line of **output** is prefixed
* `Empty` is a string used to replace blank values found in output
* `MaxWidth` is an int slice specifying the maximum width of each column.  Columns exceeding their configured width are broken at a word boundary and continued on the next line.  See below for details.

You can then pass the `Config` in using the `Format` method (signature below) to
have text formatted to your liking.

Usage
=====

    SimpleFormat(intput []string) string

    Format(input []string, config *Config) string

Using MaxWidth
==============
The MaxWidth config element allows you to configure the maximum width of each column.  An input line with a data value exceeding the maximum width for a column is formatted into multiple lines, with each long column's data broken at a word boundary and continued on the line below.  For example:

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

Columns without a MaxWidth value, or columns with a MaxWidth value of zero, expand to the maximum data width as normal.

If you want to restrict the maximum width of the entire output line, specify a negative value for exactly one of the columns in a MaxWidth specification.  The inverse of the negative value is interpreted as the desired output line width, and columnize calculates a width for the designated column so that all output lines fit in the specified width.  As a convenience, if you want the last column of data to be the one with the calculated width, specify MaxWidth as a single-element array with a negative value.

For example, if you want no output line to exceed 80 characters, you could specify:

    config := Config{MaxWidth: []int{0, -80, 15}}

This specification causes the first column to be the width of the largest piece of data in that column; the last column to be a maximum of 15 characters; and the middle column to have whatever maximum width is required so that no output line exceeds 80 characters in width.

You could also specify:

    config := Config{MaxWidth: []int{-80}}

to request that no output line exceed 80 characters, with any line breaks occurring in the last data column.


