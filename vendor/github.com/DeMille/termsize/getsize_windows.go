package termsize

import (
	"syscall"
	"unsafe"
)

//
// adapted from termbox-go:
// github.com/nsf/termbox-go
//

type (
	short int16
	word  uint16

	coord struct {
		x short
		y short
	}

	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}

	buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}
)

var (
	handle   syscall.Handle
	tmp_info buffer_info
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	proc     = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func initialize() (err error) {
	handle, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	return
}

func get_size() (w, h int, err error) {
	err = get_buffer_info(handle, &tmp_info)
	if err != nil {
		return
	}

	w = int(tmp_info.window.right - tmp_info.window.left + 1)
	h = int(tmp_info.window.bottom - tmp_info.window.top + 1)
	return
}

func get_buffer_info(h syscall.Handle, info *buffer_info) (err error) {
	retCode, _, e1 := syscall.Syscall(proc.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(retCode) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
