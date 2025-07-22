package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alperdrsnn/clime"
)

func main() {
	// Advanced Clime demonstration
	clime.Clear()

	// Custom styled header
	title := "Advanced Clime Features Showcase"
	clime.NewBanner(title, clime.BannerInfo).
		WithStyle(clime.BannerStyleDouble).
		WithColor(clime.BrightCyanColor).
		WithBorderColor(clime.BrightBlueColor).
		WithIcon("🌟").
		Println()

	fmt.Println()

	// Multi-progress bar demo
	demonstrateMultiProgressBars()

	// Complex table layouts
	demonstrateAdvancedTables()

	// Custom styled boxes with different layouts
	demonstrateAdvancedBoxes()

	// Color gradients and special effects
	demonstrateColorEffects()

	// File system simulation
	demonstrateFileSystemDemo()
}

func demonstrateMultiProgressBars() {
	clime.PrintInfoBox("Multi Progress Bars", "Demonstrating multiple concurrent progress bars")
	fmt.Println()

	multiBar := clime.NewMultiBar()

	// Add multiple progress bars
	downloadBar := clime.NewProgressBar(100).
		WithLabel("Download").
		WithStyle(clime.ProgressStyleModern).
		WithColor(clime.GreenColor).
		ShowETA(true)

	uploadBar := clime.NewProgressBar(80).
		WithLabel("Upload  ").
		WithStyle(clime.ProgressStyleArrow).
		WithColor(clime.BlueColor).
		ShowRate(true)

	processBar := clime.NewProgressBar(60).
		WithLabel("Process ").
		WithStyle(clime.ProgressStyleDots).
		WithColor(clime.YellowColor)

	multiBar.AddBar(downloadBar).AddBar(uploadBar).AddBar(processBar)

	// Simulate progress
	for i := 0; i <= 100; i++ {
		if i <= 100 {
			downloadBar.Set(int64(i))
		}
		if i <= 80 {
			uploadBar.Set(int64(i))
		}
		if i <= 60 {
			processBar.Set(int64(i))
		}

		multiBar.Print()
		time.Sleep(50 * time.Millisecond)
	}

	multiBar.Println()
	clime.SuccessLine("All tasks completed!")
	fmt.Println()
}

func demonstrateAdvancedTables() {
	clime.PrintInfoBox("Advanced Tables", "Demonstrating complex table layouts and styling")
	fmt.Println()

	// Server status table
	serverTable := clime.NewTable().
		WithStyle(clime.TableStyleBold).
		WithBorderColor(clime.BlueColor).
		WithHeaderColor(clime.BrightWhiteColor).
		AddColumn("Server").
		AddColumn("Status").
		AddColumn("CPU").
		AddColumn("Memory").
		AddColumn("Uptime").
		SetColumnAlignment(2, clime.AlignRight).
		SetColumnAlignment(3, clime.AlignRight).
		SetColumnColor(1, clime.GreenColor)

	servers := [][]string{
		{"web-01", "Online", "45%", "2.1GB", "15d 4h"},
		{"web-02", "Online", "32%", "1.8GB", "15d 4h"},
		{"api-01", "Warning", "78%", "3.2GB", "12d 8h"},
		{"db-01", "Online", "23%", "4.1GB", "25d 12h"},
		{"cache-01", "Offline", "0%", "0GB", "0h"},
	}

	serverTable.AddRows(servers)
	serverTable.Println()

	// Performance metrics table with different style
	perfTable := clime.NewTable().
		WithStyle(clime.TableStyleRounded).
		WithBorderColor(clime.CyanColor).
		ShowBorders(true).
		AddColumn("Metric").
		AddColumn("Current").
		AddColumn("Target").
		AddColumn("Status").
		SetColumnColor(0, clime.BoldColor).
		SetColumnAlignment(1, clime.AlignRight).
		SetColumnAlignment(2, clime.AlignRight)

	metrics := [][]string{
		{"Response Time", "150ms", "< 200ms", "✓ Good"},
		{"Throughput", "1,250 req/s", "> 1,000 req/s", "✓ Good"},
		{"Error Rate", "0.05%", "< 0.1%", "✓ Good"},
		{"CPU Usage", "68%", "< 80%", "✓ Good"},
		{"Memory Usage", "78%", "< 85%", "⚠ Warning"},
	}

	perfTable.AddRows(metrics)
	perfTable.Println()
}

func demonstrateAdvancedBoxes() {
	clime.PrintInfoBox("Advanced Boxes", "Demonstrating complex box layouts and content organization")
	fmt.Println()

	// Multi-column layout simulation
	leftBox := clime.NewBox().
		WithTitle("📊 System Resources").
		WithWidth(35).
		WithBorderColor(clime.GreenColor).
		WithStyle(clime.BoxStyleRounded).
		AddLine("CPU Usage:    45% ████░░░░").
		AddLine("Memory:       2.1GB/8GB").
		AddLine("Disk:         125GB/500GB").
		AddLine("Network I/O:  1.2MB/s").
		AddEmptyLine().
		AddSeparator().
		AddEmptyLine().
		AddText("All systems operating within normal parameters.")

	rightBox := clime.NewBox().
		WithTitle("📈 Performance Metrics").
		WithWidth(35).
		WithBorderColor(clime.BlueColor).
		WithStyle(clime.BoxStyleRounded).
		AddLine("Requests/sec: 1,250").
		AddLine("Avg Response: 150ms").
		AddLine("Success Rate: 99.95%").
		AddLine("Active Users: 2,847").
		AddEmptyLine().
		AddSeparator().
		AddEmptyLine().
		AddText("Performance is excellent across all metrics.")

	// Print boxes side by side (simplified version)
	fmt.Println(leftBox.Render())
	fmt.Println()
	fmt.Println(rightBox.Render())
	fmt.Println()

	// Centered announcement box
	announcement := clime.NewBox().
		WithTitle("Important Announcement").
		WithAlignment(clime.BoxAlignCenter).
		WithBorderColor(clime.YellowColor).
		WithTitleColor(clime.YellowColor).
		WithStyle(clime.BoxStyleDouble).
		WithWidth(60).
		AddEmptyLine().
		AddText("Clime v1.0.0 is now available!").
		AddEmptyLine().
		AddText("New features include enhanced styling, better performance, and improved cross-platform support.").
		AddEmptyLine().
		AddText("Visit github.com/clime/clime for more information.").
		AddEmptyLine()

	announcement.Println()
}

func demonstrateColorEffects() {
	clime.PrintInfoBox("Color Effects", "Demonstrating advanced color features and effects")
	fmt.Println()

	// Rainbow text
	fmt.Println("Rainbow: " + clime.Rainbow("This text cycles through rainbow colors!"))

	// Gradient effect
	fmt.Println("Gradient: " + clime.Gradient("This text uses gradient colors", clime.RedColor, clime.BlueColor))

	// Custom RGB colors
	customColor := clime.RGB(255, 100, 150) // Custom pink
	fmt.Println("Custom RGB: " + customColor.Sprint("This uses a custom RGB color (255, 100, 150)"))

	// Hex colors
	hexColor := clime.Hex("#FF6B35") // Orange
	fmt.Println("Hex Color: " + hexColor.Sprint("This uses a hex color (#FF6B35)"))

	// Combined styles
	combined := clime.Combine(clime.Bold, clime.Underline, clime.BrightMagenta)
	fmt.Println("Combined: " + combined.Sprint("Bold + Underline + Bright Magenta"))

	fmt.Println()
}

func demonstrateFileSystemDemo() {
	clime.PrintInfoBox("File System Simulation", "Simulating a file system operation with real-time feedback")
	fmt.Println()

	// Simulate file operations
	files := []string{
		"config.json", "main.go", "utils.go", "README.md", "package.json",
		"styles.css", "index.html", "app.js", "tests/unit.go", "docs/api.md",
	}

	// Use ShowProgress helper function
	err := clime.ShowProgress(files, "Processing files", func(file string) error {
		// Simulate some processing time
		time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
		return nil
	})

	if err != nil {
		clime.ErrorLine("Error processing files: " + err.Error())
		return
	}

	fmt.Println()

	// Show completion summary
	summary := clime.NewBox().
		WithTitle("Processing Summary").
		WithBorderColor(clime.Success).
		WithTitleColor(clime.Success).
		WithStyle(clime.BoxStyleBold).
		AddLine(fmt.Sprintf("Files processed: %d", len(files))).
		AddLine("Errors: 0").
		AddLine(fmt.Sprintf("Total time: %.2fs", float64(len(files))*0.125)).
		AddEmptyLine().
		AddText("All files have been successfully processed without any errors.")

	summary.Println()

	// Final success message
	clime.SuccessLine("Advanced demo completed successfully!")
	fmt.Println()
	clime.InfoLine("This demonstrates the power and flexibility of Clime for creating sophisticated CLI interfaces.")
}
