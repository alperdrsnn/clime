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
	message          string
	bannerType       BannerType
	style            BannerStyle
	color            *Color
	borderColor      *Color
	width            int
	multiline        bool
	ResponsiveConfig *ResponsiveConfig
	useSmartSizing   bool
}

// NewBanner creates a new banner
func NewBanner(message string, bannerType BannerType) *Banner {
	banner := &Banner{
		message:        message,
		bannerType:     bannerType,
		style:          BannerStyleDefault,
		width:          SmartWidth(0.9), // Use 90% of smart width
		multiline:      true,
		useSmartSizing: true,
	}

	switch bannerType {
	case BannerSuccess:
		banner.color = Success
		banner.borderColor = Success
	case BannerWarning:
		banner.color = Warning
		banner.borderColor = Warning
	case BannerError:
		banner.color = Error
		banner.borderColor = Error
	case BannerInfo:
		banner.color = Info
		banner.borderColor = Info
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

// WithWidth sets the banner width
func (b *Banner) WithWidth(width int) *Banner {
	if width > 0 {
		b.width = width
		b.useSmartSizing = false
	}
	return b
}

// WithSmartWidth enables smart responsive width sizing
func (b *Banner) WithSmartWidth(percentage float64) *Banner {
	b.width = SmartWidth(percentage)
	b.useSmartSizing = true
	return b
}

// WithResponsiveConfig sets responsive configuration for different breakpoints
func (b *Banner) WithResponsiveConfig(config ResponsiveConfig) *Banner {
	b.ResponsiveConfig = &config
	b.useSmartSizing = true
	return b
}

// Multiline controls whether to use multiline layout for long messages
func (b *Banner) Multiline(enable bool) *Banner {
	b.multiline = enable
	return b
}

// Render renders the banner and returns the string representation
func (b *Banner) Render() string {
	if b.message == "" {
		return ""
	}

	if b.useSmartSizing {
		rm := GetResponsiveManager()
		rm.RefreshBreakpoint()
		b.calculateResponsiveSize()
	}

	b.calculateOptimalWidth()

	var result strings.Builder

	result.WriteString(b.renderTopBorder())
	result.WriteString("\n")

	lines := b.prepareLines()
	for _, line := range lines {
		result.WriteString(b.renderContentLine(line))
		result.WriteString("\n")
	}

	result.WriteString(b.renderBottomBorder())

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
			} else if getVisualWidth(currentLine.String())+1+getVisualWidth(word) <= availableWidth {
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
		if getVisualWidth(b.message) > availableWidth {
			lines = append(lines, TruncateString(b.message, availableWidth))
		} else {
			lines = append(lines, b.message)
		}
	}

	return lines
}

// calculateResponsiveSize calculates responsive banner size
func (b *Banner) calculateResponsiveSize() {
	if b.ResponsiveConfig != nil {
		rm := GetResponsiveManager()
		config := b.ResponsiveConfig.GetConfigForBreakpoint(rm.GetCurrentBreakpoint())
		if config != nil {
			if config.Width != nil {
				b.width = *config.Width
			}
			if config.Compact {
				b.multiline = false
			}
			return
		}
	}

	if b.useSmartSizing {
		b.width = SmartWidth(0.9)
	}
}

// renderTopBorder renders the top border
func (b *Banner) renderTopBorder() string {
	borderWidth := b.width - 2
	border := b.style.TopLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.TopRight

	if b.borderColor != nil {
		return b.borderColor.Sprint(border)
	}
	return border
}

// renderBottomBorder renders the bottom border
func (b *Banner) renderBottomBorder() string {
	borderWidth := b.width - 2
	border := b.style.BottomLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.BottomRight

	if b.borderColor != nil {
		return b.borderColor.Sprint(border)
	}
	return border
}

// calculateOptimalWidth calculates the optimal banner width
func (b *Banner) calculateOptimalWidth() {
	lines := b.prepareLines()
	maxLineLength := b.getMaxLineLength(lines)

	requiredWidth := maxLineLength + (2 * b.style.Padding) + 2

	if requiredWidth > b.width {
		b.width = requiredWidth
	}
}

// renderContentLine renders a single line of content with padding and border
func (b *Banner) renderContentLine(line string) string {
	availableWidth := b.width - 2

	var content strings.Builder

	if b.borderColor != nil {
		content.WriteString(b.borderColor.Sprint(b.style.Vertical))
	} else {
		content.WriteString(b.style.Vertical)
	}

	content.WriteString(strings.Repeat(" ", b.style.Padding))

	if b.color != nil {
		content.WriteString(b.color.Sprint(line))
	} else {
		content.WriteString(line)
	}

	usedWidth := (2 * b.style.Padding) + getVisualWidth(line)
	remainingSpace := availableWidth - usedWidth
	if remainingSpace > 0 {
		content.WriteString(strings.Repeat(" ", remainingSpace))
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
		if getVisualWidth(line) > maxLength {
			maxLength = getVisualWidth(line)
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
func CustomBanner(message string, textColor, borderColor *Color, style BannerStyle) *Banner {
	banner := &Banner{
		message:     message,
		bannerType:  BannerInfo,
		style:       style,
		color:       textColor,
		borderColor: borderColor,
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

	padding := (width - getVisualWidth(title) - 4) / 2
	if padding < 0 {
		padding = 0
	}

	header := strings.Repeat("=", width)
	titleLine := "=" + strings.Repeat(" ", padding) + title + strings.Repeat(" ", padding)

	for getVisualWidth(titleLine) < width-1 {
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
