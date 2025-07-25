package clime

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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

// removeANSIEscapeCodes removes ANSI escape codes from a string
func removeANSIEscapeCodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// getVisualWidth calculates the actual visual width of a string
func getVisualWidth(s string) int {
	cleanStr := removeANSIEscapeCodes(s)
	
	width := 0
	for len(cleanStr) > 0 {
		r, size := utf8.DecodeRuneInString(cleanStr)
		if r == utf8.RuneError {
			width++
		} else {
			if isWideChar(r) {
				width += 2
			} else {
				width++
			}
		}
		cleanStr = cleanStr[size:]
	}
	
	return width
}

// isWideChar checks if a Unicode character takes 2 columns in terminal
func isWideChar(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) || // Hangul Jamo
		(r >= 0x2E80 && r <= 0x2EFF) || // CJK Radicals Supplement
		(r >= 0x2F00 && r <= 0x2FDF) || // Kangxi Radicals  
		(r >= 0x2FF0 && r <= 0x2FFF) || // Ideographic Description Characters
		(r >= 0x3000 && r <= 0x303F) || // CJK Symbols and Punctuation
		(r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) || // Katakana
		(r >= 0x3100 && r <= 0x312F) || // Bopomofo
		(r >= 0x3130 && r <= 0x318F) || // Hangul Compatibility Jamo
		(r >= 0x3190 && r <= 0x319F) || // Kanbun
		(r >= 0x31A0 && r <= 0x31BF) || // Bopomofo Extended
		(r >= 0x31C0 && r <= 0x31EF) || // CJK Strokes
		(r >= 0x31F0 && r <= 0x31FF) || // Katakana Phonetic Extensions
		(r >= 0x3200 && r <= 0x32FF) || // Enclosed CJK Letters and Months
		(r >= 0x3300 && r <= 0x33FF) || // CJK Compatibility
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Unified Ideographs Extension A
		(r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0xA000 && r <= 0xA48F) || // Yi Syllables
		(r >= 0xA490 && r <= 0xA4CF) || // Yi Radicals
		(r >= 0xAC00 && r <= 0xD7AF) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE10 && r <= 0xFE1F) || // Vertical Forms
		(r >= 0xFE30 && r <= 0xFE4F) || // CJK Compatibility Forms
		(r >= 0xFE50 && r <= 0xFE6F) || // Small Form Variants
		(r >= 0xFF00 && r <= 0xFFEF) || // Halfwidth and Fullwidth Forms
		(r >= 0x1F300 && r <= 0x1F5FF) || // Miscellaneous Symbols and Pictographs (some emojis)
		(r >= 0x1F600 && r <= 0x1F64F) || // Emoticons (emojis)
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols (emojis)
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) || // Geometric Shapes Extended
		(r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Unified Ideographs Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Unified Ideographs Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Unified Ideographs Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK Unified Ideographs Extension E
		(r >= 0x2CEB0 && r <= 0x2EBEF)    // CJK Unified Ideographs Extension F
}

// PadString pads a string to the specified width using visual width calculation
func PadString(s string, width int) string {
	visualWidth := getVisualWidth(s)
	if visualWidth >= width {
		return s
	}
	padding := strings.Repeat(" ", width-visualWidth)
	return s + padding
}

// TruncateString truncates a string to the specified width with ellipsis using visual width calculation
func TruncateString(s string, width int) string {
	visualWidth := getVisualWidth(s)
	if visualWidth <= width {
		return s
	}
	if width < 3 {
		return truncateToVisualWidth(s, width)
	}
	
	truncated := truncateToVisualWidth(s, width-3)
	return truncated + "..."
}

// truncateToVisualWidth truncates string to exact visual width
func truncateToVisualWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	
	cleanStr := removeANSIEscapeCodes(s)
	currentWidth := 0
	result := ""
	
	for len(cleanStr) > 0 {
		r, size := utf8.DecodeRuneInString(cleanStr)
		charWidth := 1
		if r != utf8.RuneError && isWideChar(r) {
			charWidth = 2
		}
		
		if currentWidth + charWidth > width {
			break
		}
		
		result += string(r)
		currentWidth += charWidth
		cleanStr = cleanStr[size:]
	}
	
	return result
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
