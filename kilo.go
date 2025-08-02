package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// *** data ***//
type editorConfig struct {
	screenrows   int
	screencols   int
	orig_termios *unix.Termios
}

var E editorConfig

//*** terminal ***//

func die(err error) {
	disableRawMode()
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
	// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[2J"))
	// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[H"))
	log.Fatal(err)
}

func disableRawMode() {
	err := unix.IoctlSetTermios(unix.Stdin, unix.TCIFLUSH, E.orig_termios)
	if err != nil {
		log.Fatal(err)
	}
}

func enableRawMode() {
	var err error
	E.orig_termios, err = unix.IoctlGetTermios(unix.Stdin, unix.TIOCGETA)
	if err != nil {
		die(err)
	}

	raw := E.orig_termios

	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN
	raw.Iflag &^= unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP
	raw.Cflag |= unix.CS8
	raw.Oflag &^= unix.OPOST
	// raw.Cc[unix.VMIN] = 0
	// raw.Cc[unix.VTIME] = 1

	err = unix.IoctlSetTermios(unix.Stdin, unix.TIOCSETA, raw)
	if err != nil {
		log.Fatal(err)
	}
}

func editorReadKey() int {
	var buffer [1]byte
	var cc int
	var err error
	for cc, err = os.Stdin.Read(buffer[:]); cc != 1; cc, err = os.Stdin.Read(buffer[:]) {
	}
	if err != nil {
		die(err)
	}

	return int(buffer[0])
}

func getWindowSize(rows *int, cols *int) int {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil || ws.Col == 0 {
		return -1
	}
	*cols = int(ws.Col)
	*rows = int(ws.Row)
	return 0
}

func editorProcessKeyPress() {
	c := editorReadKey()
	switch c {
	case ('q' & 0x1f):
		fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
		// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[2J"))
		// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[H"))
		os.Exit(0)
	}
}

// *** output ***//

func editorDrawRows() {
	var y int
	for y = 0; y < E.screenrows; y++ {
		syscall.Write(int(os.Stdout.Fd()), []byte("~\r\n"))
	}
}

func editorRefreshScreen() {
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")
	// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[2J"))
	// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[H"))

	editorDrawRows()

	fmt.Fprint(os.Stdout, "\x1b[H")
	// syscall.Write(int(os.Stdout.Fd()), []byte("\x1b[H"))
}

//*** input ***//
//*** init ***//

func initEditor() {
	if getWindowSize(&E.screenrows, &E.screencols) == -1 {
		die(errors.New("getWindowSIze"))
	}
}

func main() {
	enableRawMode()

	initEditor()
	for {
		editorRefreshScreen()
		editorProcessKeyPress()
	}

}
