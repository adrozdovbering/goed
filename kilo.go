package main

import (
	"os"
)

func main() {
	var c [1]byte
	var n int

	for {
		n, _ = os.Stdin.Read(c[:])
		if c[0] == 'q' {
			break
		}
		if n != 1 {
			continue
		}

	}
}
