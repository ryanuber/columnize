# Termsize

Termsize is a cross platform go library to get terminal size.

It is adapted from the syscalls implemented in [termbox-go](https://github.com/nsf/termbox-go), but doesn't require a full blown termbox setup.

### Usage

Install with `go get -u github.com/demille/termsize`

```go
package main

import (
    "fmt"
    "github.com/demille/termsize"
)

func main() {
    if err := termsize.Init(); err != nil {
        panic(err)
    }

    w, h, err := termsize.Size()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Size: %d X %d \n", w, h)
    // Size: 110 X 30
}
```

Boom, terminal size. There isn't anything else to this package.

### Reference

[godoc.org/github.com/DeMille/termsize](https://godoc.org/github.com/DeMille/termsize)

### License

The MIT License (MIT)

Copyright (c) 2015 Sterling DeMille &lt;sterlingdemille@gmail.com&gt;

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.