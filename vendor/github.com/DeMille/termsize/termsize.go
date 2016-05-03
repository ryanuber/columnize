/*
Package termsize is a cross platform go library to get terminal size.

It is adapted from the syscalls implemented in termbox-go, but doesn't require a full blown termbox setup.

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
*/
package termsize

//
import (
	"errors"
)

// Tracks initialization status
var IsInit bool

// Initializes the termsize package.
// Needs to be called before requesting size.
func Init() (err error) {
	err = initialize()
	if err != nil {
		return
	}

	IsInit = true
	return
}

// Returns the terminal size
func Size() (w, h int, err error) {
	if !IsInit {
		err = errors.New("termsize not yet iniitialied")
		return
	}

	return get_size()
}
