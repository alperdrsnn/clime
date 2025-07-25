package clime

import (
	"fmt"
	"strings"
)

type TableStyle struct {
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
	TableStyleDefault = TableStyle{
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
	TableStyleRounded = TableStyle{
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
	TableStyleBold = TableStyle{
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
	TableStyleDouble = TableStyle{
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
	TableStyleSimple = TableStyle{
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
	TableStyleMinimal = TableStyle{
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

type TableAlignment int

const (
	AlignLeft TableAlignment = iota
	AlignCenter
	AlignRight
)

type TableColumn struct {
	Header    string
	Width     int
	Alignment TableAlignment
	Color     *Color
}

type Table struct {
	columns          []TableColumn
	rows             [][]string
	style            TableStyle
	headerColor      *Color
	borderColor      *Color
	showHeader       bool
	showBorders      bool
	padding          int
	autoResize       bool
	maxWidth         int
	ResponsiveConfig *ResponsiveConfig
	useSmartSizing   bool
}

// NewTable creates a new table
func NewTable() *Table {
	return &Table{
		columns:        make([]TableColumn, 0),
		rows:           make([][]string, 0),
		style:          TableStyleDefault,
		headerColor:    BoldColor,
		borderColor:    DimColor,
		showHeader:     true,
		showBorders:    true,
		padding:        SmartPadding(),
		autoResize:     true,
		maxWidth:       SmartWidth(0.95), // Use 95% of smart width
		useSmartSizing: true,
	}
}

// WithStyle sets the table style
func (t *Table) WithStyle(style TableStyle) *Table {
	t.style = style
	return t
}

// WithHeaderColor sets the header text color
func (t *Table) WithHeaderColor(color *Color) *Table {
	t.headerColor = color
	return t
}

// WithBorderColor sets the border color
func (t *Table) WithBorderColor(color *Color) *Table {
	t.borderColor = color
	return t
}

// ShowHeader controls whether to show the header row
func (t *Table) ShowHeader(show bool) *Table {
	t.showHeader = show
	return t
}

// ShowBorders controls whether to show borders
func (t *Table) ShowBorders(show bool) *Table {
	t.showBorders = show
	return t
}

// WithPadding sets the cell padding
func (t *Table) WithPadding(padding int) *Table {
	if padding >= 0 {
		t.padding = padding
	}
	return t
}

// AutoResize controls whether to auto-resize columns
func (t *Table) AutoResize(enable bool) *Table {
	t.autoResize = enable
	return t
}

// WithMaxWidth sets the maximum table width
func (t *Table) WithMaxWidth(width int) *Table {
	if width > 0 {
		t.maxWidth = width
		t.useSmartSizing = false
	}
	return t
}

// WithSmartWidth enables smart responsive width sizing
func (t *Table) WithSmartWidth(percentage float64) *Table {
	t.maxWidth = SmartWidth(percentage)
	t.useSmartSizing = true
	return t
}

// WithResponsiveConfig sets responsive configuration for different breakpoints
func (t *Table) WithResponsiveConfig(config ResponsiveConfig) *Table {
	t.ResponsiveConfig = &config
	t.useSmartSizing = true
	return t
}

// AddColumn adds a column to the table
func (t *Table) AddColumn(header string) *Table {
	t.columns = append(t.columns, TableColumn{
		Header:    header,
		Width:     0,
		Alignment: AlignLeft,
		Color:     nil,
	})
	return t
}

// AddColumnWithWidth adds a column with a specific width
func (t *Table) AddColumnWithWidth(header string, width int) *Table {
	t.columns = append(t.columns, TableColumn{
		Header:    header,
		Width:     width,
		Alignment: AlignLeft,
		Color:     nil,
	})
	return t
}

// AddColumnWithConfig adds a column with full configuration
func (t *Table) AddColumnWithConfig(column TableColumn) *Table {
	t.columns = append(t.columns, column)
	return t
}

// AddRow adds a row to the table
func (t *Table) AddRow(cells ...string) *Table {
	t.rows = append(t.rows, cells)
	return t
}

// AddRows adds multiple rows to the table
func (t *Table) AddRows(rows [][]string) *Table {
	t.rows = append(t.rows, rows...)
	return t
}

// SetColumnAlignment sets the alignment for a specific column
func (t *Table) SetColumnAlignment(columnIndex int, alignment TableAlignment) *Table {
	if columnIndex >= 0 && columnIndex < len(t.columns) {
		t.columns[columnIndex].Alignment = alignment
	}
	return t
}

// SetColumnColor sets the color for a specific column
func (t *Table) SetColumnColor(columnIndex int, color *Color) *Table {
	if columnIndex >= 0 && columnIndex < len(t.columns) {
		t.columns[columnIndex].Color = color
	}
	return t
}

// Clear clears all rows from the table
func (t *Table) Clear() *Table {
	t.rows = make([][]string, 0)
	return t
}

// Render renders the table and returns the string representation
func (t *Table) Render() string {
	if len(t.columns) == 0 {
		return ""
	}

	if t.useSmartSizing {
		rm := GetResponsiveManager()
		rm.RefreshBreakpoint()
		t.calculateResponsiveSize()
	}

	t.calculateColumnWidths()

	var result strings.Builder

	if t.showBorders {
		result.WriteString(t.renderTopBorder())
		result.WriteString("\n")
	}

	if t.showHeader {
		result.WriteString(t.renderHeaderRow())
		result.WriteString("\n")

		if t.showBorders {
			result.WriteString(t.renderHeaderSeparator())
			result.WriteString("\n")
		}
	}

	for i, row := range t.rows {
		result.WriteString(t.renderDataRow(row))
		result.WriteString("\n")

		if t.showBorders && i < len(t.rows)-1 {
			//@TODO: Add row separators
		}
	}

	if t.showBorders {
		result.WriteString(t.renderBottomBorder())
	}

	return result.String()
}

// Print renders and prints the table
func (t *Table) Print() {
	fmt.Print(t.Render())
}

// Println renders and prints the table with a newline
func (t *Table) Println() {
	fmt.Println(t.Render())
}

// calculateColumnWidths calculates optimal column widths
func (t *Table) calculateColumnWidths() {
	if !t.autoResize {
		return
	}

	for i, column := range t.columns {
		if column.Width == 0 {
			t.columns[i].Width = getVisualWidth(column.Header)
		}
	}

	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(t.columns) && getVisualWidth(cell) > t.columns[i].Width {
				t.columns[i].Width = getVisualWidth(cell)
			}
		}
	}

	for i := range t.columns {
		t.columns[i].Width += t.padding * 2
	}

	totalWidth := t.calculateTotalWidth()
	if totalWidth > t.maxWidth {
		t.adjustColumnWidths(totalWidth)
	}
}

// calculateTotalWidth calculates the total table width
func (t *Table) calculateTotalWidth() int {
	totalWidth := 0
	for _, column := range t.columns {
		totalWidth += column.Width
	}

	if t.showBorders {
		totalWidth += len(t.columns) + 1
	}

	return totalWidth
}

// calculateResponsiveSize calculates responsive table size
func (t *Table) calculateResponsiveSize() {
	if t.ResponsiveConfig != nil {
		rm := GetResponsiveManager()
		config := t.ResponsiveConfig.GetConfigForBreakpoint(rm.GetCurrentBreakpoint())
		if config != nil {
			if config.Width != nil {
				t.maxWidth = *config.Width
			}
			if config.Padding != nil {
				t.padding = *config.Padding
			}
			if config.Compact {
				t.padding = min(t.padding, 1)
				t.showBorders = false
			}
			return
		}
	}

	if t.useSmartSizing {
		t.maxWidth = SmartWidth(0.95)
		t.padding = SmartPadding()
	}
}

// adjustColumnWidths adjusts column widths to fit within maxWidth
func (t *Table) adjustColumnWidths(totalWidth int) {
	if len(t.columns) == 0 {
		return
	}

	excess := totalWidth - t.maxWidth
	perColumn := excess / len(t.columns)

	for i := range t.columns {
		t.columns[i].Width -= perColumn
		if t.columns[i].Width < 3 {
			t.columns[i].Width = 3
		}
	}
}

// renderTopBorder renders the top border of the table
func (t *Table) renderTopBorder() string {
	if len(t.columns) == 0 {
		return ""
	}

	var border strings.Builder
	border.WriteString(t.style.TopLeft)

	for i, column := range t.columns {
		border.WriteString(strings.Repeat(t.style.Horizontal, column.Width))
		if i < len(t.columns)-1 {
			border.WriteString(t.style.TopTee)
		}
	}

	border.WriteString(t.style.TopRight)

	if t.borderColor != nil {
		return t.borderColor.Sprint(border.String())
	}
	return border.String()
}

// renderBottomBorder renders the bottom border of the table
func (t *Table) renderBottomBorder() string {
	if len(t.columns) == 0 {
		return ""
	}

	var border strings.Builder
	border.WriteString(t.style.BottomLeft)

	for i, column := range t.columns {
		border.WriteString(strings.Repeat(t.style.Horizontal, column.Width))
		if i < len(t.columns)-1 {
			border.WriteString(t.style.BottomTee)
		}
	}

	border.WriteString(t.style.BottomRight)

	if t.borderColor != nil {
		return t.borderColor.Sprint(border.String())
	}
	return border.String()
}

// renderHeaderSeparator renders the separator between header and data
func (t *Table) renderHeaderSeparator() string {
	if len(t.columns) == 0 {
		return ""
	}

	var border strings.Builder
	border.WriteString(t.style.LeftTee)

	for i, column := range t.columns {
		border.WriteString(strings.Repeat(t.style.Horizontal, column.Width))
		if i < len(t.columns)-1 {
			border.WriteString(t.style.Cross)
		}
	}

	border.WriteString(t.style.RightTee)

	if t.borderColor != nil {
		return t.borderColor.Sprint(border.String())
	}
	return border.String()
}

// renderHeaderRow renders the header row
func (t *Table) renderHeaderRow() string {
	var row strings.Builder

	if t.showBorders {
		if t.borderColor != nil {
			row.WriteString(t.borderColor.Sprint(t.style.Vertical))
		} else {
			row.WriteString(t.style.Vertical)
		}
	}

	for _, column := range t.columns {
		cell := t.formatCell(column.Header, column.Width, column.Alignment)
		if t.headerColor != nil {
			cell = t.headerColor.Sprint(cell)
		}
		row.WriteString(cell)

		if t.showBorders {
			if t.borderColor != nil {
				row.WriteString(t.borderColor.Sprint(t.style.Vertical))
			} else {
				row.WriteString(t.style.Vertical)
			}
		}
	}

	return row.String()
}

// renderDataRow renders a data row
func (t *Table) renderDataRow(rowData []string) string {
	var row strings.Builder

	if t.showBorders {
		if t.borderColor != nil {
			row.WriteString(t.borderColor.Sprint(t.style.Vertical))
		} else {
			row.WriteString(t.style.Vertical)
		}
	}

	for i, column := range t.columns {
		cellData := ""
		if i < len(rowData) {
			cellData = rowData[i]
		}

		cell := t.formatCell(cellData, column.Width, column.Alignment)
		if column.Color != nil {
			cell = column.Color.Sprint(cell)
		}
		row.WriteString(cell)

		if t.showBorders {
			if t.borderColor != nil {
				row.WriteString(t.borderColor.Sprint(t.style.Vertical))
			} else {
				row.WriteString(t.style.Vertical)
			}
		}
	}

	return row.String()
}

// formatCell formats a cell with proper alignment and padding
func (t *Table) formatCell(content string, width int, alignment TableAlignment) string {
	if getVisualWidth(content) > width-t.padding*2 {
		content = TruncateString(content, width-t.padding*2)
	}

	contentWidth := getVisualWidth(content)
	totalPadding := width - contentWidth
	leftPadding := t.padding
	rightPadding := totalPadding - leftPadding

	switch alignment {
	case AlignCenter:
		leftPadding = totalPadding / 2
		rightPadding = totalPadding - leftPadding
	case AlignRight:
		leftPadding = totalPadding - t.padding
		rightPadding = t.padding
	}

	return strings.Repeat(" ", leftPadding) + content + strings.Repeat(" ", rightPadding)
}

// SimpleTable creates a simple table from headers and rows
func SimpleTable(headers []string, rows [][]string) string {
	table := NewTable()

	for _, header := range headers {
		table.AddColumn(header)
	}

	table.AddRows(rows)

	return table.Render()
}

// PrintSimpleTable prints a simple table
func PrintSimpleTable(headers []string, rows [][]string) {
	fmt.Print(SimpleTable(headers, rows))
}

// KeyValueTable creates a two-column key-value table
func KeyValueTable(data map[string]string) string {
	table := NewTable().
		AddColumn("Key").
		AddColumn("Value").
		SetColumnColor(0, BoldColor)

	for key, value := range data {
		table.AddRow(key, value)
	}

	return table.Render()
}

// PrintKeyValueTable prints a key-value table
func PrintKeyValueTable(data map[string]string) {
	fmt.Print(KeyValueTable(data))
}
