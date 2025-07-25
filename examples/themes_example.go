package main

import (
	"fmt"
	"github.com/alperdrsnn/clime"
)

func main() {
	clime.Header("Clime Theme System Demo")

	ShowAllThemes()
	demonstrateThemeSwitching()
	demonstrateThemedComponents()

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}

func ShowAllThemes() {
	clime.InfoBanner("Available Themes")
	fmt.Println()

	clime.ShowAllThemes()
}

func demonstrateThemeSwitching() {
	clime.InfoBanner("Theme Switching Demo")
	fmt.Println()

	themes := []string{"dark", "light", "colorful", "minimal", "ocean"}

	for _, themeName := range themes {
		fmt.Printf("Switching to %s theme...\n", clime.BoldColor.Sprint(themeName))

		if err := clime.SetTheme(themeName); err != nil {
			fmt.Printf("Error switching to theme %s: %s\n", clime.BoldColor.Sprint(themeName), err.Error())
			continue
		}

		clime.SuccessLine(fmt.Sprintf("Success message with %s theme", themeName))
		clime.WarningLine(fmt.Sprintf("Warning message with %s theme", themeName))
		clime.ErrorLine(fmt.Sprintf("Error message with %s theme", themeName))
		clime.InfoLine(fmt.Sprintf("Info message with %s theme", themeName))
		fmt.Println()
	}

	clime.SetTheme("dark")
}

func demonstrateThemedComponents() {
	clime.InfoBanner("Themed Components Demo")
	fmt.Println()

	themes := []string{"colorful", "ocean", "minimal"}

	for _, themeName := range themes {
		clime.SetTheme(themeName)
		fmt.Printf("=== %s Theme\n", clime.BoldColor.Sprint(themeName))

		clime.ThemedBanner("This is a themed success banner!", clime.BannerSuccess).Println()
		clime.ThemedBanner("This is a themed warning banner!", clime.BannerWarning).Println()
		clime.ThemedBanner("This is a themed error banner!", clime.BannerError).Println()
		clime.ThemedBanner("This is a themed info banner!", clime.BannerInfo).Println()

		theme := clime.GetTheme()
		clime.NewBox().
			WithTitle(fmt.Sprintf("%s Theme Box", theme.Name)).
			WithBorderColor(theme.Border).
			AddLine("This box uses the current theme colors").
			AddLine(fmt.Sprintf("Primary: %s", theme.Primary.Sprint("Sample Text"))).
			AddLine(fmt.Sprintf("Secondary: %s", theme.Secondary.Sprint("Sample Text"))).
			AddLine(fmt.Sprintf("Success: %s", theme.Success.Sprint("Sample Text"))).
			Println()

		clime.NewTable().
			AddColumn("Component").
			AddColumn("Status").
			AddColumn("Theme").
			AddRow("Banner", theme.Success.Sprint("Active"), theme.Name).
			AddRow("Box", theme.Info.Sprint("Ready"), theme.Name).
			AddRow("Table", theme.Warning.Sprint("Demo"), theme.Name).
			Println()

		fmt.Println()
	}

	demonstrateCustomTheming()
}

func demonstrateCustomTheming() {
	clime.SetTheme("dark")

	clime.InfoBanner("Custom Theme Colors")
	fmt.Println()

	customRed := clime.RGB(255, 100, 100)
	customGreen := clime.RGB(100, 255, 100)
	customBlue := clime.RGB(100, 100, 255)

	fmt.Printf("Custom Red: %s\n", customRed.Sprint("Sample Text"))
	fmt.Printf("Custom Green: %s\n", customGreen.Sprint("Sample Text"))
	fmt.Printf("Custom Blue: %s\n", customBlue.Sprint("Sample Text"))
	fmt.Println()

	hexOrange := clime.Hex("#FF8C00")
	hexPurple := clime.Hex("#9932CC")
	hexTeal := clime.Hex("#008B8B")

	fmt.Printf("Hex Orange: %s\n", hexOrange.Sprint("Sample Text"))
	fmt.Printf("Hex Purple: %s\n", hexPurple.Sprint("Sample Text"))
	fmt.Printf("Hex Teal: %s\n", hexTeal.Sprint("Sample Text"))
	fmt.Println()

	fmt.Printf("Gradient effect: %s\n", clime.Gradient("Hello Gradient World!", customRed, customBlue))
	fmt.Printf("Rainbow effect: %s\n", clime.Rainbow("Rainbow Colors Are Amazing!"))
	fmt.Println()
}
