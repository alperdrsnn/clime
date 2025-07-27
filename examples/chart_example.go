package main

import (
	"fmt"
	"github.com/alperdrsnn/clime"
	"math"
	"math/rand"
	"time"
)

func main() {
	clime.Header("Clime Charts Demo")
	clime.SetTheme("colorful")

	demonstrateBarCharts()
	demonstratePieCharts()
	demonstrateHistograms()

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

func demonstrateHistograms() {
	clime.InfoBanner("Histograms")
	fmt.Println()

	data1 := generateNormalDistribution(100, 50, 10)
	data2 := generateUniformDistribution(80, 10, 90)

	clime.NewHistogram("Normal Distribution", data1).
		WithBins(12).
		WithWidth(70).
		WithColor(clime.BrightBlueColor).
		Println()

	clime.NewHistogram("Uniform Distribution", data2).
		WithBins(10).
		WithWidth(70).
		WithColor(clime.BrightGreenColor).
		Println()

	responseTimes := []float64{
		12.3, 15.7, 18.9, 22.1, 25.4, 28.7, 31.2, 34.8, 37.5, 41.2,
		44.6, 48.3, 52.1, 55.8, 59.4, 63.7, 67.9, 71.2, 75.6, 79.3,
		83.1, 87.4, 91.8, 95.2, 98.7, 102.4, 106.1, 109.8, 113.5, 117.2,
	}

	clime.NewHistogram("API Response Times (ms)", responseTimes).
		WithBins(8).
		WithWidth(60).
		WithColor(clime.BrightYellowColor).
		Println()
}

func generateNormalDistribution(count int, mean, stddev float64) []float64 {
	rand.Seed(time.Now().UnixNano())
	data := make([]float64, count)

	for i := 0; i < count; i++ {
		u1 := rand.Float64()
		u2 := rand.Float64()
		z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		data[i] = mean + z0*stddev
	}

	return data
}

func generateUniformDistribution(count int, min, max float64) []float64 {
	rand.Seed(time.Now().UnixNano() + 1)
	data := make([]float64, count)

	for i := 0; i < count; i++ {
		data[i] = min + rand.Float64()*(max-min)
	}

	return data
}
