//go:build windows

package clime

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	// enable ANSI control sequences for Windows terminal
	h := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	err := windows.GetConsoleMode(h, &originalMode)
	if err != nil {
		return
	}
	err = windows.SetConsoleMode(h,
		originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if err != nil {
		return
	}
}
