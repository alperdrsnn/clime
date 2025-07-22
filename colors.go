package clime

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
)

const (
	Reset = "\033[0m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Strike    = "\033[9m"
)

type ColorFunc func(string) string

type Color struct {
	code     string
	disabled bool
}

// NewColor creates a new color with the given ANSI code
func NewColor(code string) *Color {
	return &Color{
		code:     code,
		disabled: !term.IsTerminal(int(os.Stdout.Fd())),
	}
}

// Sprint applies the color to a string and returns it
func (c *Color) Sprint(s string) string {
	if c.disabled {
		return s
	}
	return c.code + s + Reset
}

// Sprintf applies the color to a formatted string
func (c *Color) Sprintf(format string, args ...interface{}) string {
	return c.Sprint(fmt.Sprintf(format, args...))
}

// Print prints the colored string to stdout
func (c *Color) Print(s string) {
	fmt.Print(c.Sprint(s))
}

// Printf prints the formatted colored string to stdout
func (c *Color) Printf(format string, args ...interface{}) {
	fmt.Print(c.Sprintf(format, args...))
}

// Println prints the colored string with a newline
func (c *Color) Println(s string) {
	fmt.Println(c.Sprint(s))
}

// Disable disables color output for this color
func (c *Color) Disable() *Color {
	c.disabled = true
	return c
}

// Enable enables color output for this color
func (c *Color) Enable() *Color {
	c.disabled = false
	return c
}

// IsDisabled returns true if color is disabled
func (c *Color) IsDisabled() bool {
	return c.disabled
}

var (
	BlackColor   = NewColor(Black)
	RedColor     = NewColor(Red)
	GreenColor   = NewColor(Green)
	YellowColor  = NewColor(Yellow)
	BlueColor    = NewColor(Blue)
	MagentaColor = NewColor(Magenta)
	CyanColor    = NewColor(Cyan)
	WhiteColor   = NewColor(White)

	BrightBlackColor   = NewColor(BrightBlack)
	BrightRedColor     = NewColor(BrightRed)
	BrightGreenColor   = NewColor(BrightGreen)
	BrightYellowColor  = NewColor(BrightYellow)
	BrightBlueColor    = NewColor(BrightBlue)
	BrightMagentaColor = NewColor(BrightMagenta)
	BrightCyanColor    = NewColor(BrightCyan)
	BrightWhiteColor   = NewColor(BrightWhite)

	BoldColor      = NewColor(Bold)
	DimColor       = NewColor(Dim)
	ItalicColor    = NewColor(Italic)
	UnderlineColor = NewColor(Underline)
	BlinkColor     = NewColor(Blink)
	ReverseColor   = NewColor(Reverse)
	StrikeColor    = NewColor(Strike)
)

var (
	Success = GreenColor
	Warning = YellowColor
	Error   = RedColor
	Info    = BlueColor
	Muted   = DimColor
)

// RGB creates a color from RGB values (0-255)
func RGB(r, g, b int) *Color {
	code := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return NewColor(code)
}

// Hex creates a color from a hex string (e.g., "#FF0000" or "FF0000")
func Hex(hex string) *Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return NewColor("")
	}

	r, err1 := strconv.ParseInt(hex[0:2], 16, 64)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return NewColor("")
	}

	return RGB(int(r), int(g), int(b))
}

// Combine combines multiple color codes
func Combine(codes ...string) *Color {
	combined := strings.Join(codes, "")
	return NewColor(combined)
}

// DisableColors globally disables color output
func DisableColors() {
	colors := []*Color{
		BlackColor, RedColor, GreenColor, YellowColor, BlueColor, MagentaColor, CyanColor, WhiteColor,
		BrightBlackColor, BrightRedColor, BrightGreenColor, BrightYellowColor, BrightBlueColor,
		BrightMagentaColor, BrightCyanColor, BrightWhiteColor,
		BoldColor, DimColor, ItalicColor, UnderlineColor, BlinkColor, ReverseColor, StrikeColor,
	}
	for _, color := range colors {
		color.Disable()
	}
}

// EnableColors globally enables color output
func EnableColors() {
	colors := []*Color{
		BlackColor, RedColor, GreenColor, YellowColor, BlueColor, MagentaColor, CyanColor, WhiteColor,
		BrightBlackColor, BrightRedColor, BrightGreenColor, BrightYellowColor, BrightBlueColor,
		BrightMagentaColor, BrightCyanColor, BrightWhiteColor,
		BoldColor, DimColor, ItalicColor, UnderlineColor, BlinkColor, ReverseColor, StrikeColor,
	}
	for _, color := range colors {
		color.Enable()
	}
}

// Gradient creates a gradient effect across text
func Gradient(text string, startColor, endColor *Color) string {
	if len(text) == 0 {
		return ""
	}

	var result strings.Builder
	for i, char := range text {
		if i%2 == 0 {
			result.WriteString(startColor.Sprint(string(char)))
		} else {
			result.WriteString(endColor.Sprint(string(char)))
		}
	}
	return result.String()
}

// Rainbow applies rainbow colors to text
func Rainbow(text string) string {
	colors := []*Color{RedColor, YellowColor, GreenColor, CyanColor, BlueColor, MagentaColor}
	var result strings.Builder
	for i, char := range text {
		color := colors[i%len(colors)]
		result.WriteString(color.Sprint(string(char)))
	}
	return result.String()
}
