package clime

import (
	"fmt"
	"strings"
)

type BoxStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	Cross       string
	TopTee      string
	BottomTee   string
	LeftTee     string
	RightTee    string
}

var (
	BoxStyleDefault = BoxStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
		Cross:       "┼",
		TopTee:      "┬",
		BottomTee:   "┴",
		LeftTee:     "├",
		RightTee:    "┤",
	}
	BoxStyleRounded = BoxStyle{
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Horizontal:  "─",
		Vertical:    "│",
		Cross:       "┼",
		TopTee:      "┬",
		BottomTee:   "┴",
		LeftTee:     "├",
		RightTee:    "┤",
	}
	BoxStyleBold = BoxStyle{
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
		Horizontal:  "━",
		Vertical:    "┃",
		Cross:       "╋",
		TopTee:      "┳",
		BottomTee:   "┻",
		LeftTee:     "┣",
		RightTee:    "┫",
	}
	BoxStyleDouble = BoxStyle{
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		Vertical:    "║",
		Cross:       "╬",
		TopTee:      "╦",
		BottomTee:   "╩",
		LeftTee:     "╠",
		RightTee:    "╣",
	}
	BoxStyleSimple = BoxStyle{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Horizontal:  "-",
		Vertical:    "|",
		Cross:       "+",
		TopTee:      "+",
		BottomTee:   "+",
		LeftTee:     "+",
		RightTee:    "+",
	}
	BoxStyleMinimal = BoxStyle{
		TopLeft:     " ",
		TopRight:    " ",
		BottomLeft:  " ",
		BottomRight: " ",
		Horizontal:  " ",
		Vertical:    " ",
		Cross:       " ",
		TopTee:      " ",
		BottomTee:   " ",
		LeftTee:     " ",
		RightTee:    " ",
	}
)

type BoxAlignment int

const (
	BoxAlignLeft BoxAlignment = iota
	BoxAlignCenter
	BoxAlignRight
)

type Box struct {
	content     []string
	title       string
	style       BoxStyle
	alignment   BoxAlignment
	padding     int
	width       int
	height      int
	color       *Color
	borderColor *Color
	titleColor  *Color
	autoSize    bool
	showBorder  bool
}

// NewBox creates a new box
func NewBox() *Box {
	terminal := NewTerminal()
	return &Box{
		content:     make([]string, 0),
		style:       BoxStyleDefault,
		alignment:   BoxAlignLeft,
		padding:     1,
		width:       terminal.Width() - 4, // Leave margins
		color:       nil,
		borderColor: DimColor,
		titleColor:  BoldColor,
		autoSize:    true,
		showBorder:  true,
	}
}

// WithTitle sets the box title
func (b *Box) WithTitle(title string) *Box {
	b.title = title
	return b
}

// WithStyle sets the box style
func (b *Box) WithStyle(style BoxStyle) *Box {
	b.style = style
	return b
}

// WithAlignment sets the text alignment
func (b *Box) WithAlignment(alignment BoxAlignment) *Box {
	b.alignment = alignment
	return b
}

// WithPadding sets the internal padding
func (b *Box) WithPadding(padding int) *Box {
	if padding >= 0 {
		b.padding = padding
	}
	return b
}

// WithWidth sets the box width
func (b *Box) WithWidth(width int) *Box {
	if width > 0 {
		b.width = width
		b.autoSize = false
	}
	return b
}

// WithHeight sets the box height
func (b *Box) WithHeight(height int) *Box {
	if height > 0 {
		b.height = height
		b.autoSize = false
	}
	return b
}

// WithColor sets the text color
func (b *Box) WithColor(color *Color) *Box {
	b.color = color
	return b
}

// WithBorderColor sets the border color
func (b *Box) WithBorderColor(color *Color) *Box {
	b.borderColor = color
	return b
}

// WithTitleColor sets the title color
func (b *Box) WithTitleColor(color *Color) *Box {
	b.titleColor = color
	return b
}

// AutoSize controls whether to auto-size the box
func (b *Box) AutoSize(enable bool) *Box {
	b.autoSize = enable
	return b
}

// ShowBorder controls whether to show the border
func (b *Box) ShowBorder(show bool) *Box {
	b.showBorder = show
	return b
}

// AddLine adds a single line of content
func (b *Box) AddLine(line string) *Box {
	b.content = append(b.content, line)
	return b
}

// AddLines adds multiple lines of content
func (b *Box) AddLines(lines ...string) *Box {
	b.content = append(b.content, lines...)
	return b
}

// AddText adds text content, automatically wrapping long lines
func (b *Box) AddText(text string) *Box {
	if text == "" {
		b.content = append(b.content, "")
		return b
	}

	availableWidth := b.width - (b.padding * 2)
	if b.showBorder {
		availableWidth -= 2
	}

	if availableWidth <= 0 {
		availableWidth = 20
	}

	lines := wrapText(text, availableWidth)
	b.content = append(b.content, lines...)
	return b
}

// AddEmptyLine adds an empty line
func (b *Box) AddEmptyLine() *Box {
	b.content = append(b.content, "")
	return b
}

// AddSeparator adds a horizontal separator line
func (b *Box) AddSeparator() *Box {
	availableWidth := b.width - (b.padding * 2)
	if b.showBorder {
		availableWidth -= 2
	}

	separator := strings.Repeat("─", availableWidth)
	if b.borderColor != nil {
		separator = b.borderColor.Sprint(separator)
	}

	b.content = append(b.content, separator)
	return b
}

// Clear clears all content
func (b *Box) Clear() *Box {
	b.content = make([]string, 0)
	return b
}

// Render renders the box and returns the string representation
func (b *Box) Render() string {
	if b.autoSize {
		b.calculateSize()
	}

	var result strings.Builder

	if b.showBorder {
		result.WriteString(b.renderTopBorder())
		result.WriteString("\n")
	}

	contentLines := b.prepareContentLines()
	for _, line := range contentLines {
		result.WriteString(b.renderContentLine(line))
		result.WriteString("\n")
	}

	if b.showBorder {
		result.WriteString(b.renderBottomBorder())
	}

	return result.String()
}

// Print renders and prints the box
func (b *Box) Print() {
	fmt.Print(b.Render())
}

// Println renders and prints the box with a newline
func (b *Box) Println() {
	fmt.Println(b.Render())
}

// calculateSize automatically calculates the optimal box size
func (b *Box) calculateSize() {
	if len(b.content) == 0 {
		b.width = 20
		b.height = 3
		return
	}

	maxLineLength := 0
	for _, line := range b.content {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}

	requiredWidth := maxLineLength + (b.padding * 2)
	if b.showBorder {
		requiredWidth += 2
	}

	if b.title != "" && len(b.title)+4 > requiredWidth {
		requiredWidth = len(b.title) + 4
	}

	b.width = requiredWidth

	b.height = len(b.content) + (b.padding * 2)
	if b.showBorder {
		b.height += 2
	}
}

// prepareContentLines prepares content lines for rendering
func (b *Box) prepareContentLines() []string {
	var lines []string

	for i := 0; i < b.padding; i++ {
		lines = append(lines, "")
	}

	for _, line := range b.content {
		lines = append(lines, line)
	}

	for i := 0; i < b.padding; i++ {
		lines = append(lines, "")
	}

	if !b.autoSize && b.height > 0 {
		requiredContentLines := b.height
		if b.showBorder {
			requiredContentLines -= 2
		}

		for len(lines) < requiredContentLines {
			lines = append(lines, "")
		}

		if len(lines) > requiredContentLines {
			lines = lines[:requiredContentLines]
		}
	}

	return lines
}

// renderTopBorder renders the top border with optional title
func (b *Box) renderTopBorder() string {
	borderWidth := b.width
	if b.showBorder {
		borderWidth -= 2
	}

	var border string

	if b.title != "" {
		titleLen := len(b.title)
		if titleLen+4 >= borderWidth {
			maxTitleLen := borderWidth - 4
			if maxTitleLen > 0 {
				title := TruncateString(b.title, maxTitleLen)
				border = b.style.TopLeft + "─" + title + "─" + strings.Repeat(b.style.Horizontal, borderWidth-len(title)-2) + b.style.TopRight
			} else {
				border = b.style.TopLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.TopRight
			}
		} else {
			leftPadding := (borderWidth - titleLen - 2) / 2
			rightPadding := borderWidth - titleLen - 2 - leftPadding

			border = b.style.TopLeft + strings.Repeat(b.style.Horizontal, leftPadding) + " "
			if b.titleColor != nil {
				border += b.titleColor.Sprint(b.title)
			} else {
				border += b.title
			}
			border += " " + strings.Repeat(b.style.Horizontal, rightPadding) + b.style.TopRight
		}
	} else {
		border = b.style.TopLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.TopRight
	}

	if b.borderColor != nil {
		if b.title != "" && b.titleColor != nil {
			return b.borderColor.Sprint(border)
		}
		return b.borderColor.Sprint(border)
	}

	return border
}

// renderBottomBorder renders the bottom border
func (b *Box) renderBottomBorder() string {
	borderWidth := b.width
	if b.showBorder {
		borderWidth -= 2
	}

	border := b.style.BottomLeft + strings.Repeat(b.style.Horizontal, borderWidth) + b.style.BottomRight

	if b.borderColor != nil {
		return b.borderColor.Sprint(border)
	}
	return border
}

// renderContentLine renders a single content line
func (b *Box) renderContentLine(line string) string {
	availableWidth := b.width
	if b.showBorder {
		availableWidth -= 2
	}

	if len(line) > availableWidth {
		line = TruncateString(line, availableWidth)
	}

	alignedLine := b.alignText(line, availableWidth)

	if b.color != nil {
		alignedLine = b.color.Sprint(alignedLine)
	}

	var result string
	if b.showBorder {
		leftBorder := b.style.Vertical
		rightBorder := b.style.Vertical

		if b.borderColor != nil {
			leftBorder = b.borderColor.Sprint(leftBorder)
			rightBorder = b.borderColor.Sprint(rightBorder)
		}

		result = leftBorder + alignedLine + rightBorder
	} else {
		result = alignedLine
	}

	return result
}

// alignText aligns text within the specified width
func (b *Box) alignText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}

	padding := width - textLen

	switch b.alignment {
	case BoxAlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
	case BoxAlignRight:
		return strings.Repeat(" ", padding) + text
	default:
		return text + strings.Repeat(" ", padding)
	}
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
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

	return lines
}

// SimpleBox creates a simple box with content
func SimpleBox(title, content string) string {
	return NewBox().
		WithTitle(title).
		AddText(content).
		Render()
}

// InfoBox creates an info-styled box
func InfoBox(title, content string) string {
	return NewBox().
		WithTitle(title).
		WithBorderColor(Info).
		WithTitleColor(Info).
		AddText(content).
		Render()
}

// WarningBox creates a warning-styled box
func WarningBox(title, content string) string {
	return NewBox().
		WithTitle(title).
		WithBorderColor(Warning).
		WithTitleColor(Warning).
		AddText(content).
		Render()
}

// ErrorBox creates an error-styled box
func ErrorBox(title, content string) string {
	return NewBox().
		WithTitle(title).
		WithBorderColor(Error).
		WithTitleColor(Error).
		AddText(content).
		Render()
}

// SuccessBox creates a success-styled box
func SuccessBox(title, content string) string {
	return NewBox().
		WithTitle(title).
		WithBorderColor(Success).
		WithTitleColor(Success).
		AddText(content).
		Render()
}

// PrintSimpleBox prints a simple box
func PrintSimpleBox(title, content string) {
	fmt.Print(SimpleBox(title, content))
}

// PrintInfoBox prints an info box
func PrintInfoBox(title, content string) {
	fmt.Print(InfoBox(title, content))
}

// PrintWarningBox prints a warning box
func PrintWarningBox(title, content string) {
	fmt.Print(WarningBox(title, content))
}

// PrintErrorBox prints an error box
func PrintErrorBox(title, content string) {
	fmt.Print(ErrorBox(title, content))
}

// PrintSuccessBox prints a success box
func PrintSuccessBox(title, content string) {
	fmt.Print(SuccessBox(title, content))
}
