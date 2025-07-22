package clime

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

type ProgressBarStyle struct {
	LeftBorder  string
	RightBorder string
	Filled      string
	Empty       string
	Pointer     string
}

var (
	ProgressStyleDefault = ProgressBarStyle{
		LeftBorder:  "[",
		RightBorder: "]",
		Filled:      "█",
		Empty:       "░",
		Pointer:     "",
	}
	ProgressStyleModern = ProgressBarStyle{
		LeftBorder:  "▐",
		RightBorder: "▌",
		Filled:      "▓",
		Empty:       "░",
		Pointer:     "",
	}
	ProgressStyleArrow = ProgressBarStyle{
		LeftBorder:  "(",
		RightBorder: ")",
		Filled:      "=",
		Empty:       "-",
		Pointer:     ">",
	}
	ProgressStyleDots = ProgressBarStyle{
		LeftBorder:  "[",
		RightBorder: "]",
		Filled:      "●",
		Empty:       "○",
		Pointer:     "",
	}
	ProgressStyleBlock = ProgressBarStyle{
		LeftBorder:  "▕",
		RightBorder: "▏",
		Filled:      "▉",
		Empty:       " ",
		Pointer:     "",
	}
	ProgressStyleGradient = ProgressBarStyle{
		LeftBorder:  "[",
		RightBorder: "]",
		Filled:      "█",
		Empty:       "▁",
		Pointer:     "",
	}
)

type ProgressBar struct {
	current     int64
	total       int64
	width       int
	style       ProgressBarStyle
	color       *Color
	bgColor     *Color
	label       string
	showPercent bool
	showCount   bool
	showRate    bool
	showETA     bool
	startTime   time.Time
	mu          sync.RWMutex
	finished    bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64) *ProgressBar {
	terminal := NewTerminal()
	width := terminal.Width() - 30
	if width < 20 {
		width = 20
	}

	return &ProgressBar{
		total:       total,
		width:       width,
		style:       ProgressStyleDefault,
		color:       GreenColor,
		bgColor:     DimColor,
		showPercent: true,
		showCount:   true,
		startTime:   time.Now(),
	}
}

// WithWidth sets the progress bar width
func (p *ProgressBar) WithWidth(width int) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	if width > 0 {
		p.width = width
	}
	return p
}

// WithStyle sets the progress bar style
func (p *ProgressBar) WithStyle(style ProgressBarStyle) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.style = style
	return p
}

// WithColor sets the progress bar color
func (p *ProgressBar) WithColor(color *Color) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.color = color
	return p
}

// WithBackgroundColor sets the background color
func (p *ProgressBar) WithBackgroundColor(color *Color) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bgColor = color
	return p
}

// WithLabel sets a label for the progress bar
func (p *ProgressBar) WithLabel(label string) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.label = label
	return p
}

// ShowPercent controls whether to show percentage
func (p *ProgressBar) ShowPercent(show bool) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.showPercent = show
	return p
}

// ShowCount controls whether to show count (current/total)
func (p *ProgressBar) ShowCount(show bool) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.showCount = show
	return p
}

// ShowRate controls whether to show rate (items/sec)
func (p *ProgressBar) ShowRate(show bool) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.showRate = show
	return p
}

// ShowETA controls whether to show estimated time of arrival
func (p *ProgressBar) ShowETA(show bool) *ProgressBar {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.showETA = show
	return p
}

// Set sets the current progress value
func (p *ProgressBar) Set(current int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if current > p.total {
		current = p.total
	}
	if current < 0 {
		current = 0
	}
	p.current = current
	p.finished = current >= p.total
}

// Add increments the current progress by the given amount
func (p *ProgressBar) Add(delta int64) {
	p.mu.RLock()
	current := p.current
	p.mu.RUnlock()
	p.Set(current + delta)
}

// Increment increments the current progress by 1
func (p *ProgressBar) Increment() {
	p.Add(1)
}

// Render renders the progress bar and returns the string representation
func (p *ProgressBar) Render() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var progress float64
	if p.total > 0 {
		progress = float64(p.current) / float64(p.total)
	}
	if progress > 1.0 {
		progress = 1.0
	}

	var parts []string

	if p.label != "" {
		parts = append(parts, p.label)
	}

	bar := p.buildBar(progress)
	parts = append(parts, bar)

	if p.showPercent {
		percentage := fmt.Sprintf("%3.0f%%", progress*100)
		parts = append(parts, percentage)
	}

	if p.showCount {
		count := fmt.Sprintf("(%d/%d)", p.current, p.total)
		parts = append(parts, count)
	}

	if p.showRate {
		elapsed := time.Since(p.startTime).Seconds()
		if elapsed > 0 {
			rate := float64(p.current) / elapsed
			rateStr := fmt.Sprintf("%.1f/s", rate)
			parts = append(parts, rateStr)
		}
	}

	if p.showETA && !p.finished {
		eta := p.calculateETA()
		if eta > 0 {
			etaStr := p.formatDuration(eta)
			parts = append(parts, "ETA "+etaStr)
		}
	}

	return strings.Join(parts, " ")
}

// Print renders and prints the progress bar
func (p *ProgressBar) Print() {
	rendered := p.Render()
	if p.IsFinished() {
		fmt.Print("\r" + rendered + "\n")
	} else {
		fmt.Print("\r" + rendered)
	}
}

// Println renders and prints the progress bar with a newline
func (p *ProgressBar) Println() {
	fmt.Println(p.Render())
}

// Finish completes the progress bar
func (p *ProgressBar) Finish() {
	p.Set(p.total)
	fmt.Print("\r" + p.Render() + "\n")
}

// IsFinished returns true if the progress bar is finished
func (p *ProgressBar) IsFinished() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.finished
}

// GetCurrent returns the current progress value
func (p *ProgressBar) GetCurrent() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.current
}

// GetTotal returns the total progress value
func (p *ProgressBar) GetTotal() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.total
}

// SetTotal sets a new total value
func (p *ProgressBar) SetTotal(total int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.total = total
	if p.current > p.total {
		p.current = p.total
	}
	p.finished = p.current >= p.total
}

// buildBar builds the visual progress bar
func (p *ProgressBar) buildBar(progress float64) string {
	filledLength := int(math.Round(float64(p.width) * progress))
	emptyLength := p.width - filledLength

	var filled string
	if filledLength > 0 {
		filled = strings.Repeat(p.style.Filled, filledLength-len(p.style.Pointer))
		if p.style.Pointer != "" && progress > 0 && progress < 1.0 {
			filled += p.style.Pointer
		} else if filledLength > 0 {
			filled += strings.Repeat(p.style.Filled, len(p.style.Pointer))
		}
	}

	empty := strings.Repeat(p.style.Empty, emptyLength)

	if p.color != nil {
		filled = p.color.Sprint(filled)
	}
	if p.bgColor != nil {
		empty = p.bgColor.Sprint(empty)
	}

	return p.style.LeftBorder + filled + empty + p.style.RightBorder
}

// calculateETA calculates estimated time of arrival
func (p *ProgressBar) calculateETA() time.Duration {
	if p.current == 0 {
		return 0
	}

	elapsed := time.Since(p.startTime)
	remaining := p.total - p.current
	rate := float64(p.current) / elapsed.Seconds()

	if rate <= 0 {
		return 0
	}

	eta := time.Duration(float64(remaining)/rate) * time.Second
	return eta
}

// formatDuration formats a duration for display
func (p *ProgressBar) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// MultiBar represents multiple progress bars
type MultiBar struct {
	bars []*ProgressBar
	mu   sync.RWMutex
}

// NewMultiBar creates a new multi-progress bar
func NewMultiBar() *MultiBar {
	return &MultiBar{
		bars: make([]*ProgressBar, 0),
	}
}

// AddBar adds a progress bar to the multi-bar
func (m *MultiBar) AddBar(bar *ProgressBar) *MultiBar {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bars = append(m.bars, bar)
	return m
}

// Render renders all progress bars
func (m *MultiBar) Render() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lines []string
	for _, bar := range m.bars {
		lines = append(lines, bar.Render())
	}
	return strings.Join(lines, "\n")
}

// Print renders and prints all progress bars
func (m *MultiBar) Print() {
	output := m.Render()
	lines := strings.Count(output, "\n") + 1

	if lines > 1 {
		MoveCursorUp(lines - 1)
	}

	fmt.Print("\r" + output)
}

// Println renders and prints all progress bars with a final newline
func (m *MultiBar) Println() {
	fmt.Println(m.Render())
}

// ShowProgress shows a progress bar for a slice operation
func ShowProgress[T any](items []T, label string, fn func(T) error) error {
	bar := NewProgressBar(int64(len(items))).WithLabel(label)

	for _, item := range items {
		err := fn(item)
		if err != nil {
			bar.Println()
			return err
		}
		bar.Increment()
		bar.Print()
	}

	bar.Finish()
	return nil
}

// ShowProgressWithStyle shows a progress bar with custom style
func ShowProgressWithStyle[T any](items []T, label string, style ProgressBarStyle, fn func(T) error) error {
	bar := NewProgressBar(int64(len(items))).WithLabel(label).WithStyle(style)

	for _, item := range items {
		err := fn(item)
		if err != nil {
			bar.Println()
			return err
		}
		bar.Increment()
		bar.Print()
	}

	bar.Finish()
	return nil
}
