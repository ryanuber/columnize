Columnize
=========

Easy column-formatted output for golang

[![Build Status](https://travis-ci.org/ryanuber/columnize.png)](https://travis-ci.org/ryanuber/columnize)

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

You can fine-tune the format of the output by calling the `Format` method. This
lets you set spacing and delimiter selection.

Usage
=====

```go
SimpleFormat(intput interface{}) (string, error)

Format(input interface{}, delim string, space string) (string, error)
```
