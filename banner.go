package clime

import (
	"fmt"
	"strings"
)

type BannerStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	Padding     int
}

var (
	BannerStyleDefault = BannerStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
		Padding:     1,
	}
	BannerStyleRounded = BannerStyle{
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Horizontal:  "─",
		Vertical:    "│",
		Padding:     1,
	}
	BannerStyleBold = BannerStyle{
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
		Horizontal:  "━",
		Vertical:    "┃",
		Padding:     1,
	}
	BannerStyleDouble = BannerStyle{
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		Vertical:    "║",
		Padding:     1,
	}
	BannerStyleSimple = BannerStyle{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		Vertical:    "|",
		Padding:     1,
	}
)

type BannerType int

const (
	BannerSuccess BannerType = iota
	BannerWarning
	BannerError
	BannerInfo
)

type Banner struct {
	message     string
	bannerType  BannerType
	style       BannerStyle
	color       *Color
	borderColor *Color
	icon        string
	width       int
	multiline   bool
}

// NewBanner creates a new banner
func NewBanner(message string, bannerType BannerType) *Banner {
	terminal := NewTerminal()

	banner := &Banner{
		message:    message,
		bannerType: bannerType,
		style:      BannerStyleDefault,
		width:      terminal.Width() - 4, // Leave some margin
		multiline:  true,
	}

	switch bannerType {
	case BannerSuccess:
		banner.color = Success
		banner.borderColor = Success
		banner.icon = "✓"
	case BannerWarning:
		banner.color = Warning
		banner.borderColor = Warning
		banner.icon = "⚠"
	case BannerError:
		banner.color = Error
		banner.borderColor = Error
		banner.icon = "✗"
	case BannerInfo:
		banner.color = Info
		banner.borderColor = Info
		banner.icon = "ℹ"
	}

	return banner
}

// WithStyle sets the banner style
func (b *Banner) WithStyle(style BannerStyle) *Banner {
	b.style = style
	return b
}

// WithColor sets the text color
func (b *Banner) WithColor(color *Color) *Banner {
	b.color = color
	return b
}

// WithBorderColor sets the border color
func (b *Banner) WithBorderColor(color *Color) *Banner {
	b.borderColor = color
	return b
}

// WithIcon sets a custom icon
func (b *Banner) WithIcon(icon string) *Banner {
	b.icon = icon
	return b
}

// WithWidth sets the banner width
func (b *Banner) WithWidth(width int) *Banner {
	if width > 0 {
		b.width = width
	}
	return b
}

// Multiline controls whether to use multiline layout for long messages
func (b *Banner) Multiline(enable bool) *Banner {
	b.multiline = enable
	return b
}

// Render renders the banner and returns the string representation
func (b *Banner) Render() string {
	lines := b.prepareLines()
	if len(lines) == 0 {
		return ""
	}

	var result strings.Builder

	// Top border
	result.WriteString(b.renderTopBorder(lines))
	result.WriteString("\n")

	// Content lines
	for _, line := range lines {
		result.WriteString(b.renderContentLine(line))
		result.WriteString("\n")
	}

	// Bottom border
	result.WriteString(b.renderBottomBorder(lines))

	return result.String()
}

// Print renders and prints the banner
func (b *Banner) Print() {
	fmt.Print(b.Render())
}

// Println renders and prints the banner with a newline
func (b *Banner) Println() {
	fmt.Println(b.Render())
}

// prepareLines prepares the message lines for rendering
func (b *Banner) prepareLines() []string {
	if b.message == "" {
		return []string{}
	}

	// Calculate available width for content
	availableWidth := b.width - (2 * b.style.Padding) - 2 // 2 for borders

	if b.icon != "" {
		availableWidth -= len(b.icon) + 1 // +1 for space after icon
	}

	if availableWidth <= 0 {
		availableWidth = 10
	}

	var lines []string

	if b.multiline {
		words := strings.Fields(b.message)
		var currentLine strings.Builder

		for _, word := range words {
			if currentLine.Len() == 0 {
				currentLine.WriteString(word)
			} else if currentLine.Len()+1+len(word) <= availableWidth {
				currentLine.WriteString(" " + word)
			} else {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
			}
		}

		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
	} else {
		if len(b.message) > availableWidth {
			lines = append(lines, TruncateString(b.message, availableWidth))
		} else {
			lines = append(lines, b.message)
		}
	}

	return lines
}

// renderTopBorder renders the top border
func (b *Banner) renderTopBorder(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	maxLineLength := b.getMaxLineLength(lines)
	borderWidth := maxLineLength + (2 * b.style.Padding)

	if b.icon != "" {
		borderWidth += len(b.icon) + 1
	}

	border := b.style.TopLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.TopRight

	if b.borderColor != nil {
		return b.borderColor.Sprint(border)
	}
	return border
}

// renderBottomBorder renders the bottom border
func (b *Banner) renderBottomBorder(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	maxLineLength := b.getMaxLineLength(lines)
	borderWidth := maxLineLength + (2 * b.style.Padding)

	if b.icon != "" {
		borderWidth += len(b.icon) + 1
	}

	border := b.style.BottomLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.BottomRight

	if b.borderColor != nil {
		return b.borderColor.Sprint(border)
	}
	return border
}

// renderContentLine renders a content line
func (b *Banner) renderContentLine(line string) string {
	var content strings.Builder

	if b.borderColor != nil {
		content.WriteString(b.borderColor.Sprint(b.style.Vertical))
	} else {
		content.WriteString(b.style.Vertical)
	}

	content.WriteString(strings.Repeat(" ", b.style.Padding))

	if b.icon != "" {
		if b.color != nil {
			content.WriteString(b.color.Sprint(b.icon) + " ")
		} else {
			content.WriteString(b.icon + " ")
		}
	}

	if b.color != nil {
		content.WriteString(b.color.Sprint(line))
	} else {
		content.WriteString(line)
	}

	iconWidth := 0
	if b.icon != "" {
		iconWidth = len(b.icon) + 1
	}

	paddingNeeded := b.width - len(line) - iconWidth - (2 * b.style.Padding) - 2
	if paddingNeeded > 0 {
		content.WriteString(strings.Repeat(" ", paddingNeeded))
	}

	content.WriteString(strings.Repeat(" ", b.style.Padding))

	if b.borderColor != nil {
		content.WriteString(b.borderColor.Sprint(b.style.Vertical))
	} else {
		content.WriteString(b.style.Vertical)
	}

	return content.String()
}

// getMaxLineLength gets the maximum length among all lines
func (b *Banner) getMaxLineLength(lines []string) int {
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	return maxLength
}

// SuccessBanner creates and displays a success banner
func SuccessBanner(message string) {
	NewBanner(message, BannerSuccess).Println()
}

// WarningBanner creates and displays a warning banner
func WarningBanner(message string) {
	NewBanner(message, BannerWarning).Println()
}

// ErrorBanner creates and displays an error banner
func ErrorBanner(message string) {
	NewBanner(message, BannerError).Println()
}

// InfoBanner creates and displays an info banner
func InfoBanner(message string) {
	NewBanner(message, BannerInfo).Println()
}

// SuccessLine prints a simple success message with icon
func SuccessLine(message string) {
	fmt.Println(Success.Sprint("✓ " + message))
}

// WarningLine prints a simple warning message with icon
func WarningLine(message string) {
	fmt.Println(Warning.Sprint("⚠ " + message))
}

// ErrorLine prints a simple error message with icon
func ErrorLine(message string) {
	fmt.Println(Error.Sprint("✗ " + message))
}

// InfoLine prints a simple info message with icon
func InfoLine(message string) {
	fmt.Println(Info.Sprint("ℹ " + message))
}

// CustomBanner creates a custom banner with specific colors and style
func CustomBanner(message, icon string, textColor, borderColor *Color, style BannerStyle) *Banner {
	banner := &Banner{
		message:     message,
		bannerType:  BannerInfo,
		style:       style,
		color:       textColor,
		borderColor: borderColor,
		icon:        icon,
		width:       NewTerminal().Width() - 4,
		multiline:   true,
	}

	return banner
}

// Header creates a header-style banner
func Header(title string) {
	terminal := NewTerminal()
	width := terminal.Width()
	if width > 80 {
		width = 80
	}

	padding := (width - len(title) - 4) / 2
	if padding < 0 {
		padding = 0
	}

	header := strings.Repeat("=", width)
	titleLine := "=" + strings.Repeat(" ", padding) + title + strings.Repeat(" ", padding)

	for len(titleLine) < width-1 {
		titleLine += " "
	}
	titleLine += "="

	fmt.Println(BoldColor.Sprint(header))
	fmt.Println(BoldColor.Sprint(titleLine))
	fmt.Println(BoldColor.Sprint(header))
}

// Separator prints a simple separator line
func Separator() {
	terminal := NewTerminal()
	width := terminal.Width()
	if width > 80 {
		width = 80
	}
	fmt.Println(DimColor.Sprint(strings.Repeat("─", width)))
}
