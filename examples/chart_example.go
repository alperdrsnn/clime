package main

import (
	"fmt"
	"github.com/alperdrsnn/clime"
)

func main() {
	clime.Header("Clime Charts Demo")
	clime.SetTheme("colorful")

	demonstrateBarCharts()
	demonstratePieCharts()

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}

func demonstrateBarCharts() {
	clime.InfoBanner("Bar Charts")
	fmt.Println()

	clime.NewBarChart("Sales by Region").
		AddData("North", 85.5, clime.BlueColor).
		AddData("South", 72.3, clime.GreenColor).
		AddData("East", 91.2, clime.YellowColor).
		AddData("West", 68.7, clime.RedColor).
		AddData("Central", 79.4, clime.MagentaColor).
		SetHorizontal(true).
		SetShowValues(true).
		Println()

	clime.NewBarChart("Monthly Performance").
		AddData("Jan", 45, nil).
		AddData("Feb", 67, nil).
		AddData("Mar", 89, nil).
		AddData("Apr", 76, nil).
		AddData("May", 92, nil).
		AddData("Jun", 83, nil).
		WithHeight(12).
		SetHorizontal(false).
		SetShowValues(true).
		Println()

	clime.NewBarChart("Server Response Times (ms)").
		AddData("API-1", 45.2, clime.BrightGreenColor).
		AddData("API-2", 78.9, clime.YellowColor).
		AddData("API-3", 123.4, clime.BrightRedColor).
		AddData("API-4", 67.1, clime.BrightBlueColor).
		WithWidth(60).
		SetHorizontal(true).
		SetShowValues(true).
		Println()
}

func demonstratePieCharts() {
	clime.InfoBanner("Pie Chart")
	fmt.Println()

	clime.NewPieChart("Browser Usage").
		AddData("Chrome", 48.5, clime.BlueColor).
		AddData("Test", 20.0, clime.MagentaColor).
		AddData("Firefox", 18.2, clime.RedColor).
		AddData("Safari", 9.8, clime.GreenColor).
		AddData("Edge", 3.5, clime.YellowColor).
		WithRadius(6).
		SetShowPercentages(true).
		SetShowLegend(true).
		Println()
}
