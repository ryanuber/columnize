Columnize
=========

Easy column-formatted output for golang

[![Build Status](https://travis-ci.org/ryanuber/columnize.png)](https://travis-ci.org/ryanuber/columnize)
[![Coverage Status](https://coveralls.io/repos/ryanuber/columnize/badge.png?branch=master)](https://coveralls.io/r/ryanuber/columnize?branch=master)

Columnize is a really small Go package that makes building CLI's a little bit
easier. In some CLI designs, you want to output a number similar items in a
human-readable way with nicely aligned columns. However, figuring out how wide
to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

```go
package main

import (
    "fmt"
    "github.com/ryanuber/columnize"
)

func main() {
    output := []string{
        "Name | Gender | Age",
        "Bob | Male | 38",
        "Sally | Female | 26",
    }
    result, _ := columnize.SimpleFormat(output)
    fmt.Println(result)
}
```

As you can see, you just give it a list of strings and a delimiter.
And the result:

```
Name   Gender  Age
Bob    Male    38
Sally  Female  26
```

Columnize is tolerant of missing or empty fields, or even empty lines, so
passing in extra lines for spacing should show up as you would expect.

Columnize will also accept a plain string as input. This makes it easy if you
already have a CLI that can build up some output and just pass it through
Columnize, like so

```go
output := "Name | Gender | Age\n"
output += "Bob | Male | 38\n"
output += "Sally | Female | 26\n"

result, _ := columnize.SimpleFormat(output)
fmt.Println(result)
```

# Configuration

Columnize is configured using a `Config`, which can be obtained by calling the
`DefaultConfig()` method. You can then tweak the settings in the resulting
`Config`:

```
config := columnize.DefaultConfig()
config.Delim = "|"
config.Glue = "  "
```

You can then pass the `Config` in using the `Format` method (signature below) to
have text formatted to your liking.

Usage
=====

```go
SimpleFormat(intput interface{}) (string, error)

Format(input interface{}, config *Config) (string, error)
```
