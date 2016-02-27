// +build !windows

package termsize

import (
	"os"
	"syscall"
	"unsafe"
)

//
// adapted from termbox-go:
// github.com/nsf/termbox-go
//

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

var f *os.File

func initialize() (err error) {
	f, err = os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	return err
}

func get_size() (int, int, error) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 0, 0, errno
	}

	return int(ws.Col), int(ws.Row), nil
}
