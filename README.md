columnize
=========

Columnize is a really small Go package that makes building CLI's a little bit
easier. In some CLI designs, you want to output a number similar items in a
human-readable way with nicely aligned columns. However, figuring out how wide
to make each column is a boring problem to solve and eats your valuable time.

Here is an example:

```
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
    fmt.Println(columnize.Columnize(output, "|"))
}
```

As you can see, you just give it a list of strings and a delimiter.
And the result:

```
Name   Gender  Age
Bob    Male    38
Sally  Female  26
```
