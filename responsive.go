package clime

import (
	"sync"
)

type BreakpointSize int

const (
	BreakpointXS BreakpointSize = iota // < 60 chars
	BreakpointSM                       // 60-79 chars
	BreakpointMD                       // 80 - 119 chars
	BreakpointLG                       // 120 - 159 chars
	BreakpointXL                       // >= 160 chars
)

type Breakpoint struct {
	Size     BreakpointSize
	MinWidth int
	MaxWidth int
	Name     string
	IsActive bool
}

var (
	Breakpoints = []Breakpoint{
		{BreakpointXS, 0, 59, "xs", false},
		{BreakpointSM, 60, 79, "sm", false},
		{BreakpointMD, 80, 119, "md", false},
		{BreakpointLG, 120, 159, "lg", false},
		{BreakpointXL, 160, 999, "xl", false},
	}
)

// ResponsiveManager handles responsive behavior
type ResponsiveManager struct {
	terminal          *Terminal
	currentBreakpoint BreakpointSize
	mu                sync.RWMutex
}

var globalResponsiveManager *ResponsiveManager
var responsiveOnce sync.Once

// GetResponsiveManager returns the global responsive manager singleton
func GetResponsiveManager() *ResponsiveManager {
	responsiveOnce.Do(func() {
		globalResponsiveManager = NewResponsiveManager()
	})

	return globalResponsiveManager
}

// NewResponsiveManager creates a new responsive manager
func NewResponsiveManager() *ResponsiveManager {
	rm := &ResponsiveManager{
		terminal: NewTerminal(),
	}

	rm.updateBreakpoint()
	return rm
}

// GetCurrentBreakpoint returns the current active breakpoint
func (rm *ResponsiveManager) GetCurrentBreakpoint() BreakpointSize {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.currentBreakpoint
}

// GetCurrentBreakpointName returns the current breakpoint name
func (rm *ResponsiveManager) GetCurrentBreakpointName() string {
	bp := rm.GetCurrentBreakpoint()
	return Breakpoints[bp].Name
}

// IsBreakpoint checks if current breakpoint matches given size
func (rm *ResponsiveManager) IsBreakpoint(size BreakpointSize) bool {
	return rm.GetCurrentBreakpoint() == size
}

// IsBreakpointOrLarger checks if current breakpoint is size or larger
func (rm *ResponsiveManager) IsBreakpointOrLarger(size BreakpointSize) bool {
	return rm.GetCurrentBreakpoint() >= size
}

// IsBreakpointOrSmaller checks if current breakpoint is size or smaller
func (rm *ResponsiveManager) IsBreakpointOrSmaller(size BreakpointSize) bool {
	return rm.GetCurrentBreakpoint() <= size
}

// RefreshBreakpoint manually refreshes the current breakpoint
func (rm *ResponsiveManager) RefreshBreakpoint() {
	rm.terminal = NewTerminal()
	rm.updateBreakpoint()
}

// updateBreakpoint updates the current breakpoint based on terminal width
func (rm *ResponsiveManager) updateBreakpoint() {
	width := rm.terminal.Width()

	var newBreakpoint BreakpointSize
	for i, bp := range Breakpoints {
		if width >= bp.MinWidth && width <= bp.MaxWidth {
			newBreakpoint = BreakpointSize(i)
			break
		}
	}

	rm.mu.Lock()
	rm.currentBreakpoint = newBreakpoint

	for i := range Breakpoints {
		Breakpoints[i].IsActive = i == int(newBreakpoint)
	}
	rm.mu.Unlock()
}

// ResponsiveConfig holds responsive configuration for elements
type ResponsiveConfig struct {
	XS *ElementConfig
	SM *ElementConfig
	MD *ElementConfig
	LG *ElementConfig
	XL *ElementConfig
}

// ElementConfig defines element configuration per breakpoint
type ElementConfig struct {
	Width    *int
	Height   *int
	Padding  *int
	Margin   *int
	ShowFull bool
	Compact  bool
}

// GetConfigForBreakpoint returns the appropriate config for current breakpoint
func (rc *ResponsiveConfig) GetConfigForBreakpoint(bp BreakpointSize) *ElementConfig {
	switch bp {
	case BreakpointXS:
		if rc.XS != nil {
			return rc.XS
		}
	case BreakpointSM:
		if rc.SM != nil {
			return rc.SM
		}
		if rc.XS != nil {
			return rc.XS
		}
	case BreakpointMD:
		if rc.MD != nil {
			return rc.MD
		}

		if rc.SM != nil {
			return rc.SM
		}
		if rc.XS != nil {
			return rc.XS
		}
	case BreakpointLG:
		if rc.LG != nil {
			return rc.LG
		}

		if rc.MD != nil {
			return rc.MD
		}

		if rc.SM != nil {
			return rc.SM
		}
		if rc.XS != nil {
			return rc.XS
		}
	case BreakpointXL:
		if rc.XL != nil {
			return rc.XL
		}

		if rc.LG != nil {
			return rc.LG
		}

		if rc.MD != nil {
			return rc.MD
		}

		if rc.SM != nil {
			return rc.SM
		}
		if rc.XS != nil {
			return rc.XS
		}
	}

	return nil
}

// SmartWidth sizing functions
func SmartWidth(percentage float64) int {
	rm := GetResponsiveManager()
	terminalWidth := rm.terminal.Width()

	baseWidth := int(float64(terminalWidth) * percentage)

	switch rm.GetCurrentBreakpoint() {
	case BreakpointXS:
		return min(baseWidth, terminalWidth-2)
	case BreakpointSM:
		return min(baseWidth, terminalWidth-4)
	case BreakpointMD:
		return min(baseWidth, terminalWidth-8)
	case BreakpointLG:
		return min(baseWidth, terminalWidth-12)
	case BreakpointXL:
		maxWidth := min(terminalWidth-16, 120)
		return min(baseWidth, maxWidth)
	}

	return baseWidth
}

// SmartPadding returns appropriate padding based on screen size
func SmartPadding() int {
	rm := GetResponsiveManager()
	switch rm.GetCurrentBreakpoint() {
	case BreakpointXS:
		return 0
	case BreakpointSM:
		return 1
	case BreakpointMD:
		return 1
	case BreakpointLG:
		return 2
	case BreakpointXL:
		return 2
	}

	return 1
}

// SmartMargin returns appropriate margin based on screen size
func SmartMargin() int {
	rm := GetResponsiveManager()
	switch rm.GetCurrentBreakpoint() {
	case BreakpointXS:
		return 1
	case BreakpointSM:
		return 2
	case BreakpointMD:
		return 4
	case BreakpointLG:
		return 6
	case BreakpointXL:
		return 8
	}

	return 4
}

func IsXS() bool { return GetResponsiveManager().IsBreakpoint(BreakpointXS) }
func IsSM() bool { return GetResponsiveManager().IsBreakpoint(BreakpointSM) }
func IsMD() bool { return GetResponsiveManager().IsBreakpoint(BreakpointMD) }
func IsLG() bool { return GetResponsiveManager().IsBreakpoint(BreakpointLG) }
func IsXL() bool { return GetResponsiveManager().IsBreakpoint(BreakpointXL) }

func IsMDOrLarger() bool  { return GetResponsiveManager().IsBreakpointOrLarger(BreakpointMD) }
func IsSMOrSmaller() bool { return GetResponsiveManager().IsBreakpointOrSmaller(BreakpointSM) }

// Utility functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetOptimalColumns returns optimal number of columns for current screen size
func GetOptimalColumns(contentWidth int) int {
	rm := GetResponsiveManager()
	availableWidth := rm.terminal.Width() - SmartMargin()*2

	if contentWidth <= 0 {
		contentWidth = 20
	}

	columns := availableWidth / (contentWidth + 2)
	if columns < 1 {
		columns = 1
	}

	switch rm.GetCurrentBreakpoint() {
	case BreakpointXS:
		return min(columns, 1)
	case BreakpointSM:
		return min(columns, 2)
	case BreakpointMD:
		return min(columns, 3)
	case BreakpointLG:
		return min(columns, 4)
	case BreakpointXL:
		return min(columns, 6)
	}

	return columns
}
