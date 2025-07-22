package main

import (
	"fmt"
	"github.com/alperdrsnn/clime"
	"time"
)

func main() {
	clime.Header("Welcome to Clime Demo")
	fmt.Println()

	// Basic colors and text formatting
	fmt.Println(clime.BoldColor.Sprint("Basic Text Formatting:"))
	fmt.Println("  " + clime.Success.Sprint("Success message"))
	fmt.Println("  " + clime.Warning.Sprint("Warning message"))
	fmt.Println("  " + clime.Error.Sprint("Error message"))
	fmt.Println("  " + clime.Info.Sprint("Info message"))
	fmt.Println("  " + clime.Muted.Sprint("Muted text"))
	fmt.Println("  " + clime.BoldColor.Sprint("Bold text"))
	fmt.Println("  " + clime.UnderlineColor.Sprint("Underlined text"))
	fmt.Println("  " + clime.Rainbow("Rainbow text!"))
	fmt.Println()

	// Separators and spacing
	clime.Separator()
	fmt.Println()

	// Spinners
	fmt.Println(clime.BoldColor.Sprint("Spinners:"))

	// Basic spinner
	spinner := clime.NewSpinner().WithMessage("Loading...").Start()
	time.Sleep(2 * time.Second)
	spinner.Success("Task completed!")

	// Different spinner styles
	spinner2 := clime.NewSpinner().
		WithStyle(clime.SpinnerClock).
		WithColor(clime.MagentaColor).
		WithMessage("Processing with clock spinner...").
		Start()
	time.Sleep(1500 * time.Millisecond)
	spinner2.Success("Done with style!")

	fmt.Println()

	fmt.Println(clime.BoldColor.Sprint("Progress Bars:"))

	// Basic progress bar
	bar := clime.NewProgressBar(100).
		WithLabel("Download").
		ShowRate(true).
		ShowETA(true)

	for i := 0; i <= 100; i += 5 {
		bar.Set(int64(i))
		bar.Print()
		time.Sleep(100 * time.Millisecond)
	}
	bar.Finish()

	// Different style progress bar
	bar2 := clime.NewProgressBar(50).
		WithLabel("Upload").
		WithStyle(clime.ProgressStyleArrow).
		WithColor(clime.CyanColor)

	for i := 0; i <= 50; i += 2 {
		bar2.Set(int64(i))
		bar2.Print()
		time.Sleep(50 * time.Millisecond)
	}
	bar2.Finish()

	fmt.Println()

	// Banners
	fmt.Println(clime.BoldColor.Sprint("Banners:"))
	clime.SuccessBanner("Operation completed successfully!")
	clime.WarningBanner("This is a warning message that might be important to read.")
	clime.ErrorBanner("An error occurred during processing.")
	clime.InfoBanner("Here's some useful information for you.")
	fmt.Println()

	// Tables
	fmt.Println(clime.BoldColor.Sprint("Tables:"))
	table := clime.NewTable().
		AddColumn("Name").
		AddColumn("Status").
		AddColumn("Progress").
		SetColumnColor(1, clime.Success).
		AddRow("Task 1", "Completed", "100%").
		AddRow("Task 2", "In Progress", "75%").
		AddRow("Task 3", "Pending", "0%")

	table.Print()
	fmt.Println()

	// Boxes
	fmt.Println(clime.BoldColor.Sprint("Boxes:"))

	box := clime.NewBox().
		WithTitle("System Information").
		WithBorderColor(clime.BlueColor).
		WithStyle(clime.BoxStyleRounded).
		AddLine("OS: Go Demo").
		AddLine("Version: 1.0.0").
		AddEmptyLine().
		AddSeparator().
		AddEmptyLine().
		AddText("This is a longer text that demonstrates automatic text wrapping within boxes. The text will be wrapped to fit nicely within the box boundaries.")

	box.Println()

	clime.PrintSuccessBox("Success", "Everything is working perfectly!")
	clime.PrintWarningBox("Warning", "Please check your configuration.")
	clime.PrintErrorBox("Error", "Something went wrong.")
	fmt.Println()

	clime.SuccessLine("Demo completed! Clime makes CLI apps beautiful and easy! ✨")
}
