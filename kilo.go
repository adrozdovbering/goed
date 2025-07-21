package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

var orig_termios *unix.Termios

func disableRawMode() {
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCIFLUSH, orig_termios)
}

func enableRawMode() {
	var err error
	// read terminal attributes into termios struct.
	orig_termios, err = unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TIOCGETA)
	if err != nil {
		panic(err)
	}

	raw := orig_termios

	// switch off key printing (echo) on the screen
	raw.Lflag &^= unix.ECHO | unix.ICANON

	err = unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TIOCSETA, raw)
	if err != nil {
		panic(err)
	}
}

func main() {
	var c [1]byte
	var n int
	var err error

	enableRawMode()

	for {
		n, err = os.Stdin.Read(c[:])
		if c[0] == 'q' {
			disableRawMode()
			break
		}

		if !iscntrl(c) {
			fmt.Printf("%d\n", c)
		} else {
			fmt.Printf("%d ('%c')\n", c, c[0])
		}
		if n != 1 {
			_, err = os.Stdin.Read(c[:])
			if err != nil {
				disableRawMode()
				panic(err)
			}
			continue
		}
		if err != nil {
			disableRawMode()
			panic(err)
		}
	}
}

func iscntrl(inp [1]byte) bool {
	r := int(inp[0])
	return r >= 32 && r <= 126
}
