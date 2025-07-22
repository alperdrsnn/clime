# Clime 🎨

A beautiful and feature-rich Command Line Interface library for Go that makes building stunning terminal applications effortless.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/doc/install)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-orange.svg)]()

## ✨ Features

Clime provides a comprehensive set of tools for creating beautiful command-line interfaces:

### 🎯 Core Components
- **Terminal Utilities** - Clear screen, cursor movement, size detection
- **Rich Colors & Formatting** - 16+ colors, bold, italic, underline, rainbow effects
- **Interactive Spinners** - Multiple styles with customizable colors and messages
- **Progress Bars** - Single and multi-bar support with ETA and rate display
- **Styled Banners** - Success, warning, error, and info messages
- **Data Tables** - Formatted tables with column styling and alignment
- **Decorative Boxes** - Multiple border styles with titles and content wrapping
- **User Input** - Text, password, email, number, and confirmation prompts
- **Autocomplete** - Smart suggestions with fuzzy matching

### 🎨 Visual Elements
- **Spinner Styles**: Default, Clock, Dots, Arrow, Modern
- **Progress Styles**: Modern, Arrow, Dots, Classic
- **Box Styles**: Default, Rounded, Bold, Double
- **Banner Types**: Success, Warning, Error, Info with custom icons
- **Color Support**: Full palette including bright colors and gradients

### 🖥️ Platform Support
- **Cross-platform**: Windows, macOS, Linux
- **Terminal Detection**: Automatic TTY detection and sizing
- **Responsive Design**: Adapts to terminal width and height

## 📦 Installation

```bash
go get github.com/alperdrsnn/clime
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/alperdrsnn/clime"
)

func main() {
    // Welcome header
    clime.Header("Welcome to My App")
    
    // Colored messages
    clime.SuccessLine("Application started successfully!")
    clime.InfoLine("Loading configuration...")
    
    // Spinner with custom message
    spinner := clime.NewSpinner().
        WithMessage("Initializing...").
        WithColor(clime.BlueColor).
        Start()
    time.Sleep(2 * time.Second)
    spinner.Success("Initialization complete!")
    
    // Progress bar
    bar := clime.NewProgressBar(100).
        WithLabel("Processing").
        ShowETA(true)
        
    for i := 0; i <= 100; i += 10 {
        bar.Set(int64(i))
        bar.Print()
        time.Sleep(200 * time.Millisecond)
    }
    bar.Finish()
    
    // Interactive input
    name, _ := clime.Ask("What's your name?")
    clime.SuccessLine(fmt.Sprintf("Hello, %s! 👋", name))
}
```

## 📚 Usage Examples

### Colors and Text Formatting

```go
// Basic colors
fmt.Println(clime.Success.Sprint("Success message"))
fmt.Println(clime.Warning.Sprint("Warning message"))
fmt.Println(clime.Error.Sprint("Error message"))
fmt.Println(clime.Info.Sprint("Info message"))

// Text styling
fmt.Println(clime.BoldColor.Sprint("Bold text"))
fmt.Println(clime.UnderlineColor.Sprint("Underlined text"))
fmt.Println(clime.Rainbow("Rainbow text!"))
```

### Spinners

```go
// Basic spinner
spinner := clime.NewSpinner().
    WithMessage("Loading...").
    Start()
time.Sleep(2 * time.Second)
spinner.Success("Done!")

// Styled spinner
spinner2 := clime.NewSpinner().
    WithStyle(clime.SpinnerClock).
    WithColor(clime.MagentaColor).
    WithMessage("Processing...").
    Start()
time.Sleep(3 * time.Second)
spinner2.Success("Complete!")
```

### Progress Bars

```go
// Single progress bar
bar := clime.NewProgressBar(100).
    WithLabel("Download").
    WithStyle(clime.ProgressStyleModern).
    WithColor(clime.GreenColor).
    ShowRate(true).
    ShowETA(true)

// Multi-progress bars
multiBar := clime.NewMultiBar()
bar1 := clime.NewProgressBar(100).WithLabel("Task 1")
bar2 := clime.NewProgressBar(80).WithLabel("Task 2")
multiBar.AddBar(bar1).AddBar(bar2)
```

### Tables

```go
table := clime.NewTable().
    AddColumn("Name").
    AddColumn("Status").
    AddColumn("Progress").
    SetColumnColor(1, clime.Success).
    AddRow("Task 1", "Completed", "100%").
    AddRow("Task 2", "In Progress", "75%").
    AddRow("Task 3", "Pending", "0%")

table.Print()
```

### Boxes

```go
// Simple box
box := clime.NewBox().
    WithTitle("Information").
    WithBorderColor(clime.BlueColor).
    WithStyle(clime.BoxStyleRounded).
    AddLine("System Status: Online").
    AddLine("Version: 1.0.0").
    AddSeparator().
    AddText("This is a longer description that will be automatically wrapped.")

box.Println()

// Quick boxes
clime.PrintSuccessBox("Success", "Operation completed!")
clime.PrintWarningBox("Warning", "Check your configuration.")
clime.PrintErrorBox("Error", "Something went wrong.")
```

### Interactive Input

```go
// Basic input
name, err := clime.Ask("Enter your name:")

// Email validation
email, err := clime.AskEmail("Enter your email:")

// Password input (masked)
password, err := clime.AskPassword("Enter password:")

// Number input
age, err := clime.AskNumber("Enter your age:")

// Confirmation
confirmed, err := clime.AskConfirm("Continue?", true)

// Autocomplete
options := []string{"apple", "banana", "cherry", "date"}
choice, err := clime.AutoComplete(clime.AutoCompleteConfig{
    Label: "Choose a fruit:",
    Options: options,
    FuzzyMatch: true,
})
```

### Banners

```go
clime.SuccessBanner("Operation completed successfully!")
clime.WarningBanner("Please review the following items.")
clime.ErrorBanner("Failed to connect to the server.")
clime.InfoBanner("New update available.")

// Custom banner
clime.NewBanner("Custom Message", clime.BannerInfo).
    WithStyle(clime.BannerStyleDouble).
    WithColor(clime.CyanColor).
    WithIcon("🚀").
    Println()
```

### Terminal Utilities

```go
// Terminal information
terminal := clime.NewTerminal()
width := terminal.Width()
height := terminal.Height()
isInteractive := terminal.IsATTY()

// Screen control
clime.Clear()                    // Clear screen
clime.HideCursor()              // Hide cursor
clime.ShowCursor()              // Show cursor
clime.MoveCursorUp(3)           // Move cursor up 3 lines
clime.ClearLine()               // Clear current line
```

## 🎛️ Configuration Options

### Spinner Styles
- `SpinnerDefault` - Classic spinning animation
- `SpinnerClock` - Clock-like rotation
- `SpinnerDots` - Bouncing dots
- `SpinnerArrow` - Rotating arrow

### Progress Bar Styles
- `ProgressStyleModern` - Clean modern look
- `ProgressStyleArrow` - Arrow-based indicator
- `ProgressStyleDots` - Dotted progress
- `ProgressStyleClassic` - Traditional bar

### Box Styles
- `BoxStyleDefault` - Standard borders
- `BoxStyleRounded` - Rounded corners
- `BoxStyleBold` - Thick borders
- `BoxStyleDouble` - Double-line borders

### Banner Styles
- `BannerStyleSingle` - Single line border
- `BannerStyleDouble` - Double line border
- `BannerStyleThick` - Thick border

## 🎨 Color Palette

Clime supports a full range of colors:

**Standard Colors**: Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
**Bright Colors**: BrightRed, BrightGreen, BrightYellow, BrightBlue, etc.
**Special Effects**: Rainbow text, gradients, background colors

## Supported Features

### ✅ Current Features

- [x] **Terminal Management**: Screen clearing, cursor control, size detection
- [x] **Color System**: 16+ colors, text styling, rainbow effects  
- [x] **Spinners**: 4+ styles with customizable colors and messages
- [x] **Progress Bars**: Single and multi-bar with ETA/rate display
- [x] **Data Tables**: Column styling, alignment, color support
- [x] **Decorative Boxes**: 4+ border styles with title support
- [x] **Interactive Prompts**: Text, password, email, number input
- [x] **Autocomplete**: Fuzzy matching, customizable options
- [x] **Banners**: Success, warning, error, info with icons
- [x] **Cross-platform**: Windows, macOS, Linux support
- [x] **TTY Detection**: Automatic terminal capability detection
- [x] **Text Utilities**: Padding, truncation, wrapping
- [x] **Multiple Input Types**: Confirmation, selection, validation

### 🚧 Planned Features (TODO)

- [ ] **Enhanced Tables**: Row separators and advanced formatting
- [ ] **Windows Integration**: Native Windows API for improved terminal size detection
- [ ] **Theme System**: Predefined color themes and custom theme support
- [ ] **Advanced Layouts**: Grid systems and complex UI layouts
- [ ] **Animation System**: Custom animations and transitions
- [ ] **Configuration Files**: JSON/YAML configuration support
- [ ] **Logging Integration**: Built-in logging with styled output
- [ ] **Plugin System**: Extensible component architecture
- [ ] **Interactive Menus**: Navigation menus and selection interfaces
- [ ] **Chart Components**: Simple ASCII charts and graphs

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/alperdrsnn/clime.git
cd clime

# Install dependencies
go mod tidy

# Run examples
go run examples/basic_example.go
go run examples/interactive_example.go
go run examples/advanced_example.go
```