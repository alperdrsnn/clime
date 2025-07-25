package clime

import "fmt"

type Theme struct {
	Name       string
	Primary    *Color
	Secondary  *Color
	Success    *Color
	Warning    *Color
	Error      *Color
	Info       *Color
	Muted      *Color
	Background *Color
	Text       *Color
	Border     *Color
}

var (
	DarkTheme = &Theme{
		Name:       "Dark",
		Primary:    BrightBlueColor,
		Secondary:  BrightCyanColor,
		Success:    BrightGreenColor,
		Warning:    BrightYellowColor,
		Error:      BrightRedColor,
		Info:       BrightBlueColor,
		Muted:      DimColor,
		Background: BlackColor,
		Text:       BrightWhiteColor,
		Border:     BrightBlackColor,
	}

	LightTheme = &Theme{
		Name:       "Light",
		Primary:    BlueColor,
		Secondary:  CyanColor,
		Success:    GreenColor,
		Warning:    YellowColor,
		Error:      RedColor,
		Info:       BlueColor,
		Muted:      BlackColor,
		Background: WhiteColor,
		Text:       BlackColor,
		Border:     BlackColor,
	}

	ColorfulTheme = &Theme{
		Name:       "Colorful",
		Primary:    BrightMagentaColor,
		Secondary:  BrightCyanColor,
		Success:    BrightGreenColor,
		Warning:    BrightYellowColor,
		Error:      BrightRedColor,
		Info:       BrightBlueColor,
		Muted:      DimColor,
		Background: BlackColor,
		Text:       BrightWhiteColor,
		Border:     BrightMagentaColor,
	}

	MinimalTheme = &Theme{
		Name:       "Minimal",
		Primary:    WhiteColor,
		Secondary:  DimColor,
		Success:    WhiteColor,
		Warning:    WhiteColor,
		Error:      WhiteColor,
		Info:       WhiteColor,
		Muted:      DimColor,
		Background: BlackColor,
		Text:       WhiteColor,
		Border:     DimColor,
	}

	OceanTheme = &Theme{
		Name:       "Ocean",
		Primary:    RGB(0, 150, 255),
		Secondary:  RGB(0, 200, 200),
		Success:    RGB(0, 255, 150),
		Warning:    RGB(255, 200, 0),
		Error:      RGB(255, 100, 100),
		Info:       RGB(100, 200, 255),
		Muted:      RGB(100, 100, 150),
		Background: RGB(5, 25, 50),
		Text:       RGB(200, 230, 255),
		Border:     RGB(50, 100, 150),
	}
)

var availableThemes = map[string]*Theme{
	"dark":     DarkTheme,
	"light":    LightTheme,
	"colorful": ColorfulTheme,
	"minimal":  MinimalTheme,
	"ocean":    OceanTheme,
}

var currentTheme = DarkTheme

// SetTheme sets the active theme by name
func SetTheme(themeName string) error {
	theme, exists := availableThemes[themeName]
	if !exists {
		return fmt.Errorf("theme '%s' not found", themeName)
	}

	currentTheme = theme

	Success = theme.Success
	Warning = theme.Warning
	Error = theme.Error
	Info = theme.Info
	Muted = theme.Muted

	return nil
}

// GetTheme returns the current active theme
func GetTheme() *Theme {
	return currentTheme
}

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() []string {
	names := make([]string, 0, len(availableThemes))
	for name := range availableThemes {
		names = append(names, name)
	}
	return names
}

// ThemePreview shows a preview of a theme
func ThemePreview(themeName string) error {
	theme, exists := availableThemes[themeName]
	if !exists {
		return fmt.Errorf("theme '%s' not found", themeName)
	}

	fmt.Printf("Theme: %s\n", BoldColor.Sprint(theme.Name))
	fmt.Printf("Primary:  %s\n", theme.Primary.Sprint("Sample Text"))
	fmt.Printf("Secondary: %s\n", theme.Secondary.Sprint("Sample Text"))
	fmt.Printf("Success:  %s\n", theme.Success.Sprint("Sample Text"))
	fmt.Printf("Warning:  %s\n", theme.Warning.Sprint("Sample Text"))
	fmt.Printf("Error:    %s\n", theme.Error.Sprint("Sample Text"))
	fmt.Printf("Info:     %s\n", theme.Info.Sprint("Sample Text"))
	fmt.Printf("Muted:    %s\n", theme.Muted.Sprint("Sample Text"))
	fmt.Printf("Background: %s\n", theme.Background.Sprint("Sample Text"))
	fmt.Printf("Text:      %s\n", theme.Text.Sprint("Sample Text"))
	fmt.Printf("Border:    %s\n", theme.Border.Sprint("Sample Text"))

	return nil
}

// ShowAllThemes displays previews of all available themes
func ShowAllThemes() {
	fmt.Println(BoldColor.Sprint("Available Themes:"))
	fmt.Println()

	for _, themeName := range GetAvailableThemes() {
		ThemePreview(themeName)
		fmt.Println()
	}
}

// ThemedBanner creates a banner using current theme colors
func ThemedBanner(message string, bannerType BannerType) *Banner {
	banner := NewBanner(message, bannerType)

	switch bannerType {
	case BannerSuccess:
		banner.WithColor(currentTheme.Success).WithBorderColor(currentTheme.Border)
	case BannerWarning:
		banner.WithColor(currentTheme.Warning).WithBorderColor(currentTheme.Border)
	case BannerError:
		banner.WithColor(currentTheme.Error).WithBorderColor(currentTheme.Border)
	case BannerInfo:
		banner.WithColor(currentTheme.Info).WithBorderColor(currentTheme.Border)
	}

	return banner
}
