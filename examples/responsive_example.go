package main

import (
	"fmt"
	"time"
	"github.com/alperdrsnn/clime"
)

func main() {
	clime.Header("Clime Responsive Design Demo")
	showCurrentBreakpoint()
	
	demonstrateResponsiveBoxes()
	demonstrateResponsiveTables()
	demonstrateResponsiveBanners()
	demonstrateResponsiveProgressBars()
	
	showRefreshDemo()
}

func showCurrentBreakpoint() {
	rm := clime.GetResponsiveManager()
	
	fmt.Printf("Current breakpoint: %s\n", clime.BoldColor.Sprint(rm.GetCurrentBreakpointName()))
	fmt.Printf("Terminal type: %s\n", getTerminalType(rm.GetCurrentBreakpoint()))
	
	fmt.Println("\nBreakpoint ranges:")
	for _, bp := range clime.Breakpoints {
		status := "inactive"
		if bp.IsActive {
			status = clime.Success.Sprint("active")
		}
		fmt.Printf("  %s: %d-%d chars (%s)\n", bp.Name, bp.MinWidth, bp.MaxWidth, status)
	}
	fmt.Println()
}

func getTerminalType(bp clime.BreakpointSize) string {
	switch bp {
	case clime.BreakpointXS:
		return clime.Warning.Sprint("Very Small (Mobile/SSH)")
	case clime.BreakpointSM:
		return clime.Info.Sprint("Small Terminal")
	case clime.BreakpointMD:
		return clime.Success.Sprint("Medium Terminal")
	case clime.BreakpointLG:
		return clime.Success.Sprint("Large Terminal")
	case clime.BreakpointXL:
		return clime.Success.Sprint("Extra Large Terminal")
	default:
		return "Unknown"
	}
}

func demonstrateResponsiveBoxes() {
	clime.InfoBanner("📦 Responsive Boxes")
	fmt.Println()

	// Simple smart width box
	clime.NewBox().
		WithTitle("Smart Width Box").
		WithSmartWidth(0.8).
		AddLine("This box uses 80% of terminal width").
		AddLine("It automatically adapts to your screen size").
		AddLine("Try resizing your terminal and run again!").
		Println()

	// Multi-breakpoint configuration
	width25 := 25
	width40 := 40
	width60 := 60
	width80 := 80
	padding1 := 1
	padding2 := 2

	clime.NewBox().
		WithTitle("Multi-Breakpoint Box").
		WithResponsiveConfig(clime.ResponsiveConfig{
			XS: &clime.ElementConfig{Width: &width25, Padding: &padding1, Compact: true},
			SM: &clime.ElementConfig{Width: &width40, Padding: &padding1},
			MD: &clime.ElementConfig{Width: &width60, Padding: &padding2},
			LG: &clime.ElementConfig{Width: &width80, Padding: &padding2},
			XL: &clime.ElementConfig{Width: &width80, Padding: &padding2},
		}).
		AddLine("Width and padding change based on breakpoint").
		AddLine("XS: 25 width, 1 padding, compact mode").
		AddLine("SM: 40 width, 1 padding").
		AddLine("MD+: 60-80 width, 2 padding").
		Println()
}

func demonstrateResponsiveTables() {
	clime.InfoBanner("📊 Responsive Tables")
	fmt.Println()

	// Simple smart width table
	clime.NewTable().
		WithSmartWidth(0.9).
		AddColumn("Service").
		AddColumn("Status").
		AddColumn("CPU").
		AddRow("web-server", "Online", "45%").
		AddRow("database", "Online", "23%").
		AddRow("cache", "Warning", "78%").
		Println()

	// Compact table for small screens
	width50 := 50
	width70 := 70
	padding0 := 0
	padding1 := 1

	fmt.Println("Responsive Table (compact on small screens):")
	clime.NewTable().
		WithResponsiveConfig(clime.ResponsiveConfig{
			XS: &clime.ElementConfig{Width: &width50, Padding: &padding0, Compact: true},
			SM: &clime.ElementConfig{Width: &width70, Padding: &padding1},
			MD: &clime.ElementConfig{Width: &width70, Padding: &padding1},
		}).
		AddColumn("Task").
		AddColumn("Progress").
		AddRow("Build", "85%").
		AddRow("Test", "60%").
		AddRow("Deploy", "30%").
		Println()
}

func demonstrateResponsiveBanners() {
	clime.InfoBanner("🏷️ Responsive Banners")
	fmt.Println()

	// Smart width banner
	clime.NewBanner("This banner adapts to your terminal width automatically!", clime.BannerSuccess).
		WithSmartWidth(0.7).
		Println()

	// Responsive configuration banner
	width30 := 30
	width50 := 50
	width70 := 70

	clime.NewBanner("Multi-breakpoint banner: Different widths for different screen sizes", clime.BannerInfo).
		WithResponsiveConfig(clime.ResponsiveConfig{
			XS: &clime.ElementConfig{Width: &width30, Compact: true},
			SM: &clime.ElementConfig{Width: &width50},
			MD: &clime.ElementConfig{Width: &width70},
		}).
		Println()
}

func demonstrateResponsiveProgressBars() {
	clime.InfoBanner("📈 Responsive Progress Bars")
	fmt.Println()

	// Smart width progress bar
	fmt.Println("Smart Width Progress Bar:")
	bar := clime.NewProgressBar(100).
		WithLabel("Download").
		WithSmartWidth(0.6).
		ShowETA(true)

	for i := 0; i <= 100; i += 20 {
		bar.Set(int64(i))
		bar.Print()
		time.Sleep(200 * time.Millisecond)
	}
	bar.Finish()
	fmt.Println()

	// Responsive compact mode
	width10 := 10
	width20 := 20
	width35 := 35

	fmt.Println("Responsive Progress Bar (compact on small screens):")
	compactBar := clime.NewProgressBar(50).
		WithLabel("Upload").
		WithResponsiveConfig(clime.ResponsiveConfig{
			XS: &clime.ElementConfig{Width: &width10, Compact: true},
			SM: &clime.ElementConfig{Width: &width20, Compact: true},
			MD: &clime.ElementConfig{Width: &width35},
		})

	for i := 0; i <= 50; i += 10 {
		compactBar.Set(int64(i))
		compactBar.Print()
		time.Sleep(150 * time.Millisecond)
	}
	compactBar.Finish()
	fmt.Println()
}

func showRefreshDemo() {
	clime.InfoBanner("🔄 Manual Refresh Demo")
	fmt.Println()
	
	fmt.Println("Try resizing your terminal, then press 'r' to refresh.")
	fmt.Println("Press 'q' to quit the demo.")
	fmt.Println()
	
	for {
		var input string
		fmt.Print("Command (r=refresh, q=quit): ")
		fmt.Scanln(&input)
		
		switch input {
		case "r", "R":
			// Refresh and show updated info
			rm := clime.GetResponsiveManager()
			rm.RefreshBreakpoint()
			
			fmt.Println("\n" + clime.Success.Sprint("Breakpoint refreshed!"))
			showCurrentBreakpoint()
			
			// Show updated responsive box
			clime.NewBox().
				WithTitle("Updated Responsive Box").
				WithSmartWidth(0.7).
				AddLine(fmt.Sprintf("Current breakpoint: %s", rm.GetCurrentBreakpointName())).
				AddLine("All responsive elements updated successfully!").
				Println()

			// Updated banner
			clime.NewBanner("Responsive system is working! 🚀", clime.BannerSuccess).
				WithSmartWidth(0.8).
				Println()
			
		case "q", "Q":
			fmt.Println("Demo completed. Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Use 'r' or 'q'.")
		}
	}
}
