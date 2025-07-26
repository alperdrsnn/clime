package clime

import (
	"fmt"
	"math"
	"strings"
)

// ChartData represents data for charts
type ChartData struct {
	Label string
	Value float64
	Color *Color
}

// BarChart represents a bar chart
type BarChart struct {
	Title            string
	Data             []ChartData
	Width            int
	Height           int
	MaxValue         float64
	ShowValues       bool
	Horizontal       bool
	ResponsiveConfig *ResponsiveConfig
	useSmartSizing   bool
}

// NewBarChart creates a new bar chart
func NewBarChart(title string) *BarChart {
	return &BarChart{
		Title:          title,
		Data:           make([]ChartData, 0),
		Width:          SmartWidth(0.8),
		Height:         10,
		ShowValues:     true,
		Horizontal:     false,
		useSmartSizing: true,
	}
}

// AddData adds data to the chart
func (bc *BarChart) AddData(label string, value float64, color *Color) *BarChart {
	if color == nil {
		colors := []*Color{BlueColor, GreenColor, YellowColor, RedColor, MagentaColor, CyanColor}
		color = colors[len(bc.Data)%len(colors)]
	}

	bc.Data = append(bc.Data, ChartData{Label: label, Value: value, Color: color})

	if value > bc.MaxValue {
		bc.MaxValue = value
	}

	return bc
}

// WithWidth sets the chart width
func (bc *BarChart) WithWidth(width int) *BarChart {
	bc.Width = width
	bc.useSmartSizing = false
	return bc
}

// WithHeight sets the chart height (for horizontal charts)
func (bc *BarChart) WithHeight(height int) *BarChart {
	bc.Height = height
	return bc
}

// WithSmartWidth sets responsive width
func (bc *BarChart) WithSmartWidth(ratio float64) *BarChart {
	bc.Width = SmartWidth(ratio)
	bc.useSmartSizing = true
	return bc
}

// SetShowValues toggles value display
func (bc *BarChart) SetShowValues(show bool) *BarChart {
	bc.ShowValues = show
	return bc
}

// SetHorizontal sets chart orientation
func (bc *BarChart) SetHorizontal(horizontal bool) *BarChart {
	bc.Horizontal = horizontal
	return bc
}

// WithResponsiveConfig sets responsive configuration
func (bc *BarChart) WithResponsiveConfig(config ResponsiveConfig) *BarChart {
	bc.ResponsiveConfig = &config
	bc.useSmartSizing = true
	return bc
}

// Print renders and prints the chart
func (bc *BarChart) Print() {
	fmt.Print(bc.Render())
}

// Println renders and prints the chart with newline
func (bc *BarChart) Println() {
	fmt.Println(bc.Render())
}

// Render generates the chart string
func (bc *BarChart) Render() string {
	if len(bc.Data) == 0 {
		return ""
	}

	var result strings.Builder

	if bc.Title != "" {
		titleLine := fmt.Sprintf("%s", bc.Title)
		result.WriteString(BoldColor.Sprint(titleLine) + "\n\n")
	}

	if bc.Horizontal {
		result.WriteString(bc.renderHorizontal())
	} else {
		result.WriteString(bc.renderVertical())
	}

	return result.String()
}

// renderHorizontal renders horizontal bar chart
func (bc *BarChart) renderHorizontal() string {
	var result strings.Builder

	maxLabelWidth := 0
	for _, data := range bc.Data {
		if getVisualWidth(data.Label) > maxLabelWidth {
			maxLabelWidth = getVisualWidth(data.Label)
		}
	}

	barWidth := bc.Width - maxLabelWidth - 10
	if barWidth < 10 {
		barWidth = 10
	}

	for _, data := range bc.Data {
		label := PadString(data.Label, maxLabelWidth)
		result.WriteString(label + " ")

		percentage := data.Value / bc.MaxValue
		barLength := int(percentage * float64(barWidth))

		bar := strings.Repeat("█", barLength)
		bar += strings.Repeat("░", barWidth-barLength)

		result.WriteString(data.Color.Sprint(bar))

		if bc.ShowValues {
			valueStr := fmt.Sprintf(" %.1f", data.Value)
			result.WriteString(DimColor.Sprint(valueStr))
		}

		result.WriteString("\n")
	}

	return result.String()
}

// renderVertical renders vertical bar chart
func (bc *BarChart) renderVertical() string {
	var result strings.Builder

	barCount := len(bc.Data)
	barWidth := (bc.Width - barCount - 1) / barCount
	if barWidth < 1 {
		barWidth = 1
	}

	for row := bc.Height; row > 0; row-- {
		threshold := (float64(row) / float64(bc.Height)) * bc.MaxValue

		for i, data := range bc.Data {
			if i > 0 {
				result.WriteString(" ")
			}

			if data.Value >= threshold {
				bar := strings.Repeat("█", barWidth)
				result.WriteString(data.Color.Sprint(bar))
			} else {
				bar := strings.Repeat(" ", barWidth)
				result.WriteString(bar)
			}
		}
		result.WriteString("\n")
	}

	for i, data := range bc.Data {
		if i > 0 {
			result.WriteString(" ")
		}

		label := TruncateString(data.Label, barWidth)
		label = PadString(label, barWidth)
		result.WriteString(label)
	}
	result.WriteString("\n")

	if bc.ShowValues {
		for i, data := range bc.Data {
			if i > 0 {
				result.WriteString(" ")
			}

			valueStr := fmt.Sprintf("%.1f", data.Value)
			valueStr = TruncateString(valueStr, barWidth)
			valueStr = PadString(valueStr, barWidth)
			result.WriteString(DimColor.Sprint(valueStr))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// PieChart represents a pie chart
type PieChart struct {
	Title            string
	Data             []ChartData
	Radius           int
	ShowPercentages  bool
	ShowLegend       bool
	ResponsiveConfig *ResponsiveConfig
}

// NewPieChart creates a new pie chart
func NewPieChart(title string) *PieChart {
	return &PieChart{
		Title:           title,
		Data:            make([]ChartData, 0),
		Radius:          8,
		ShowPercentages: true,
		ShowLegend:      true,
	}
}

// AddData adds data to the pie chart
func (pc *PieChart) AddData(label string, value float64, color *Color) *PieChart {
	if color == nil {
		colors := []*Color{BlueColor, GreenColor, YellowColor, RedColor, MagentaColor, CyanColor}
		color = colors[len(pc.Data)&len(colors)]
	}

	pc.Data = append(pc.Data, ChartData{
		Label: label,
		Value: value,
		Color: color,
	})

	return pc
}

// WithRadius sets the pie chart radius
func (pc *PieChart) WithRadius(radius int) *PieChart {
	pc.Radius = radius
	return pc
}

// SetShowPercentages toggles percentage display
func (pc *PieChart) SetShowPercentages(show bool) *PieChart {
	pc.ShowPercentages = show
	return pc
}

// SetShowLegend toggles legend display
func (pc *PieChart) SetShowLegend(show bool) *PieChart {
	pc.ShowLegend = show
	return pc
}

// Print renders and prints the pie chart
func (pc *PieChart) Print() {
	fmt.Print(pc.Render())
}

// Println renders and prints the pie chart with newline
func (pc *PieChart) Println() {
	fmt.Println(pc.Render())
}

// Render generates the pie chart string
func (pc *PieChart) Render() string {
	if len(pc.Data) == 0 {
		return ""
	}

	var result strings.Builder

	if pc.Title != "" {
		titleLine := fmt.Sprintf("%s", pc.Title)
		result.WriteString(BoldColor.Sprintf(titleLine) + "\n\n")
	}

	total := 0.0
	for _, data := range pc.Data {
		total += data.Value
	}

	effectiveRadius := float64(pc.Radius)
	size := int(effectiveRadius * 2.2)

	charAspectRatio := 0.45

	for y := 0; y < size; y++ {
		line := strings.Builder{}
		for x := 0; x < int(float64(size)*2); x++ {
			fx := float64(x)/2.0 - effectiveRadius
			fy := (float64(y) - effectiveRadius) / charAspectRatio

			//distance := math.Sqrt(fx*fx + fy*fy)

			var coverage float64
			samples := 4

			for sy := 0; sy < samples; sy++ {
				for sx := 0; sx < samples; sx++ {
					sampleX := fx + (float64(sx)-1.5)/float64(samples)*0.5
					sampleY := fy + (float64(sy)-1.5)/float64(samples)*0.5/charAspectRatio

					sampleDist := math.Sqrt(sampleX*sampleX + sampleY*sampleY)
					if sampleDist <= effectiveRadius {
						coverage += 1.0 / float64(samples*samples)
					}
				}
			}
			if coverage > 0.1 {
				angle := math.Atan2(fx, -fy)
				if angle < 0 {
					angle += 2 * math.Pi
				}

				currentAngle := 0.0
				var selectedData *ChartData

				for i := range pc.Data {
					sliceAngle := (pc.Data[i].Value / total) * 2 * math.Pi
					if angle >= currentAngle && angle < currentAngle+sliceAngle {
						selectedData = &pc.Data[i]
						break
					}
					currentAngle += sliceAngle
				}

				var char string
				if coverage > 0.9 {
					char = "█"
				} else if coverage > 0.7 {
					char = "▉"
				} else if coverage > 0.5 {
					char = "▊"
				} else if coverage > 0.3 {
					char = "▋"
				} else {
					char = " "
				}

				if selectedData != nil {
					line.WriteString(selectedData.Color.Sprint(char))
				} else {
					line.WriteString(char)
				}
			} else {
				line.WriteString(" ")
			}
		}
		lineStr := strings.TrimRight(line.String(), " ")
		if lineStr != "" {
			result.WriteString(lineStr)
		}
		result.WriteString("\n")
	}

	if pc.ShowLegend {
		result.WriteString("\nLegend:\n")
		for _, data := range pc.Data {
			percentage := (data.Value / total) * 100
			legendLine := fmt.Sprintf("  %s %s", data.Color.Sprint("█"), data.Label)

			if pc.ShowPercentages {
				legendLine += fmt.Sprintf(" (%.1f%%)", percentage)
			}

			result.WriteString(legendLine + "\n")
		}
	}

	return result.String()
}
