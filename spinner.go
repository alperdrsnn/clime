package clime

import (
	"fmt"
	"sync"
	"time"
)

type SpinnerStyle struct {
	Frames   []string
	Interval time.Duration
}

var (
	SpinnerDots = SpinnerStyle{
		Frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Interval: 80 * time.Millisecond,
	}
	SpinnerLine = SpinnerStyle{
		Frames:   []string{"|", "/", "-", "\\"},
		Interval: 100 * time.Millisecond,
	}
	SpinnerArrow = SpinnerStyle{
		Frames:   []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
		Interval: 120 * time.Millisecond,
	}
	SpinnerBounce = SpinnerStyle{
		Frames:   []string{"⠁", "⠂", "⠄", "⠂"},
		Interval: 200 * time.Millisecond,
	}
	SpinnerClock = SpinnerStyle{
		Frames:   []string{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"},
		Interval: 100 * time.Millisecond,
	}
	SpinnerEarth = SpinnerStyle{
		Frames:   []string{"🌍", "🌎", "🌏"},
		Interval: 180 * time.Millisecond,
	}
	SpinnerMoon = SpinnerStyle{
		Frames:   []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
		Interval: 80 * time.Millisecond,
	}
	SpinnerRunner = SpinnerStyle{
		Frames:   []string{"🚶", "🏃"},
		Interval: 140 * time.Millisecond,
	}
	SpinnerPulse = SpinnerStyle{
		Frames:   []string{"●", "◐", "◑", "◒", "◓", "◔", "◕", "◖", "◗"},
		Interval: 100 * time.Millisecond,
	}
	SpinnerGrowVertical = SpinnerStyle{
		Frames:   []string{"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"},
		Interval: 120 * time.Millisecond,
	}
)

type Spinner struct {
	style      SpinnerStyle
	color      *Color
	message    string
	prefix     string
	suffix     string
	running    bool
	stopCh     chan bool
	mu         sync.RWMutex
	hideCursor bool
}

// NewSpinner creates a new spinner with the default style
func NewSpinner() *Spinner {
	return &Spinner{
		style:      SpinnerDots,
		color:      CyanColor,
		stopCh:     make(chan bool),
		hideCursor: true,
	}
}

// WithStyle sets the spinner style
func (s *Spinner) WithStyle(style SpinnerStyle) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.style = style
	return s
}

// WithColor sets the spinner color
func (s *Spinner) WithColor(color *Color) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.color = color
	return s
}

// WithMessage sets the spinner message
func (s *Spinner) WithMessage(message string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
	return s
}

// WithPrefix sets a prefix for the spinner
func (s *Spinner) WithPrefix(prefix string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = prefix
	return s
}

// WithSuffix sets a suffix for the spinner
func (s *Spinner) WithSuffix(suffix string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.suffix = suffix
	return s
}

// HideCursor controls whether to hide the cursor while spinning
func (s *Spinner) HideCursor(hide bool) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hideCursor = hide
	return s
}

// Start starts the spinner animation
func (s *Spinner) Start() *Spinner {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return s
	}
	s.running = true
	s.stopCh = make(chan bool)
	s.mu.Unlock()

	if s.hideCursor {
		HideCursor()
	}

	go s.animate()
	return s
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	close(s.stopCh)
	s.mu.Unlock()

	ClearLine()
	if s.hideCursor {
		ShowCursor()
	}
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Print(Success.Sprint("✓ ") + message + "\n")
}

// Error stops the spinner and shows an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Print(Error.Sprint("✗ ") + message + "\n")
}

// Warning stops the spinner and shows a warning message
func (s *Spinner) Warning(message string) {
	s.Stop()
	fmt.Print(Warning.Sprint("⚠ ") + message + "\n")
}

// Info stops the spinner and shows an info message
func (s *Spinner) Info(message string) {
	s.Stop()
	fmt.Print(Info.Sprint("ℹ ") + message + "\n")
}

// UpdateMessage updates the spinner message while it's running
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

// IsRunning returns true if the spinner is currently running
func (s *Spinner) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	ticker := time.NewTicker(s.style.Interval)
	defer ticker.Stop()

	frameIndex := 0
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.RLock()
			frame := s.style.Frames[frameIndex]
			output := s.buildOutput(frame)
			s.mu.RUnlock()

			ClearLine()
			fmt.Print(output)

			frameIndex = (frameIndex + 1) % len(s.style.Frames)
		}
	}
}

// buildOutput builds the complete spinner output string
func (s *Spinner) buildOutput(frame string) string {
	var output string

	if s.prefix != "" {
		output += s.prefix + " "
	}

	if s.color != nil {
		output += s.color.Sprint(frame)
	} else {
		output += frame
	}

	if s.message != "" {
		output += " " + s.message
	}

	if s.suffix != "" {
		output += " " + s.suffix
	}

	return output
}

// ShowSpinner shows a spinner with a message and runs the provided function
func ShowSpinner(message string, fn func() error) error {
	s := NewSpinner().WithMessage(message).Start()
	defer s.Stop()

	err := fn()
	if err != nil {
		s.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	s.Success("Done!")
	return nil
}

// ShowSpinnerWithStyle shows a spinner with custom style, message and runs the provided function
func ShowSpinnerWithStyle(style SpinnerStyle, message string, fn func() error) error {
	s := NewSpinner().WithStyle(style).WithMessage(message).Start()
	defer s.Stop()

	err := fn()
	if err != nil {
		s.Error(fmt.Sprintf("Failed: %v", err))
		return err
	}

	s.Success("Done!")
	return nil
}
