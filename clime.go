package clime

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const Version = "1.0.0"

type Terminal struct {
	width  int
	height int
	isATTY bool
}

// NewTerminal creates a new terminal instance
func NewTerminal() *Terminal {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	return &Terminal{
		width:  width,
		height: height,
		isATTY: term.IsTerminal(int(os.Stdout.Fd())),
	}
}

// Width returns the terminal width
func (t *Terminal) Width() int {
	return t.width
}

// Height returns the terminal height
func (t *Terminal) Height() int {
	return t.height
}

// IsATTY returns true if stdout is a terminal
func (t *Terminal) IsATTY() bool {
	return t.isATTY
}

// Clear clears the terminal screen
func Clear() {
	fmt.Print("\033[2J\033[H")
}

// MoveCursorUp moves the cursor up by n lines
func MoveCursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

// MoveCursorDown moves the cursor down by n lines
func MoveCursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}

// HideCursor hides the terminal cursor
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// ClearLine clears the current line
func ClearLine() {
	fmt.Print("\033[2K\r")
}

// PadString pads a string to the specified width
func PadString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := strings.Repeat(" ", width-len(s))
	return s + padding
}

// TruncateString truncates a string to the specified width with ellipsis
func TruncateString(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width < 3 {
		return s[:width]
	}
	return s[:width-3] + "..."
}

// getTerminalSize gets terminal size using syscalls for better Windows support
func getTerminalSize() (width, height int) {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			return w, h
		}
	}

	if width, height := getWindowsTerminalSize(); width > 0 && height > 0 {
		return width, height
	}

	return 80, 24
}

// getWindowsTerminalSize gets terminal size on Windows
func getWindowsTerminalSize() (width, height int) {
	cmd := exec.Command("powershell", "-Command",
		"$Host.UI.RawUI.BufferSize.Width; $Host.UI.RawUI.BufferSize.Height")

	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return 0, 0
	}

	width, err = strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, 0
	}

	height, err = strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return 0, 0
	}

	return width, height
}
