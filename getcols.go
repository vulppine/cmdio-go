// +build !windows !android !ios
// +build !wasm

package cmdio

import (
	"fmt"
	"os/exec"
	"strconv"
)

// hacky as hell, but it works better
// than the terminfo library i found
// for some specific cases

// i wonder if there's a way to actually make this
// work without having to rely on external commands
// (for Linux, at least)

// GetCols gets the number of columns in a terminal by calling
// tput to get the window size.
func GetCols() (int, error) {
	cmd := exec.Command("tput", "cols")
	col, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// convert bytes to string and exclude
	// the newline char
	s := string(col[:len(col)-1])

	// Atoi
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	// return it
	return i, nil
}

func UpOneLine() {
	fmt.Printf("\033[A")
}

func ClearLine(m string) {
	fmt.Printf("\033[K%s\n", m)
}
