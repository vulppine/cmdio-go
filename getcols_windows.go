package cmdio

import (
	"fmt"
	"os/exec"
	"strconv"
)

// see the comment on the linux version of this for more info
// on why i mildly dislike this method

// GetCols gets the current amount of columns in the command line window
// by calling PowerShell and getting host.ui.rawui.WindowSize.Width
func GetCols() (int, error) {
	cmd := exec.Command("powershell", "-command &{(get-host).ui.rawui.WindowSize.Width;}")
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

// these apparently work in windows but
// i don't have a copy of W10 installed
// to check against the latest version

func UpOneLine() {
	fmt.Printf("\033[A")
}

func ClearLine(m string) {
	fmt.Printf("\033[K%s\n", m)
}
