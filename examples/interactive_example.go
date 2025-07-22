package main

import (
	"fmt"
	"github.com/alperdrsnn/clime"
	"strings"
)

func main() {
	clime.Header("Interactive Clime Demo")
	fmt.Println()

	clime.InfoBanner("This demo will show you various interactive input components.")
	fmt.Println()

	// Basic input
	name, err := clime.Ask("What's your name?")
	if err != nil {
		clime.ErrorLine("Error reading input: " + err.Error())
		return
	}

	if name != "" {
		clime.SuccessLine(fmt.Sprintf("Hello, %s! Nice to meet you! 👋", name))
	}
	fmt.Println()

	// Required input with validation
	email, err := clime.AskEmail("Please enter your email address")
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	clime.InfoLine("Email saved: " + email)
	fmt.Println()

	// Number input
	age, err := clime.AskNumber("What's your age?")
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	clime.InfoLine(fmt.Sprintf("Age: %d", age))
	fmt.Println()

	// Password input
	password, err := clime.AskPassword("Enter a password")
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	clime.SuccessLine("Password set successfully (length: " + fmt.Sprintf("%d", len(password)) + ")")
	fmt.Println()

	// Confirmation
	confirmed, err := clime.AskConfirm("Do you want to continue?")
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}

	if !confirmed {
		clime.WarningLine("Operation cancelled by user")
		return
	}
	fmt.Println()

	// Single choice selection
	frameworks := []string{"React", "Vue", "Angular", "Svelte", "Next.js"}
	choice, err := clime.AskChoice("Which frontend framework do you prefer?", frameworks...)
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	clime.SuccessLine("You chose: " + frameworks[choice])
	fmt.Println()

	// Multi-choice selection
	languages := []string{"Go", "JavaScript", "Python", "Rust", "TypeScript", "Java", "C++"}
	choices, err := clime.AskMultiChoice("Which programming languages do you know? (select multiple)", languages...)
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}

	var selectedLanguages []string
	for _, idx := range choices {
		selectedLanguages = append(selectedLanguages, languages[idx])
	}
	clime.InfoLine("You know: " + strings.Join(selectedLanguages, ", "))
	fmt.Println()

	// Autocomplete input
	commands := []string{"init", "start", "stop", "restart", "status", "logs", "deploy", "build", "test"}
	command, err := clime.AskWithOptions("What command would you like to run?", commands)
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	if command != "" {
		clime.SuccessLine("Running command: " + command)
	}
	fmt.Println()

	// Advanced autocomplete with builder pattern
	project, err := clime.NewAutoCompleteBuilder("Enter project name").
		WithOptions([]string{"clime-demo", "web-app", "api-server", "cli-tool", "mobile-app"}).
		WithPlaceholder("e.g., my-project").
		FuzzyMatch(true).
		Required(true).
		Ask()
	if err != nil {
		clime.ErrorLine("Error: " + err.Error())
		return
	}
	clime.InfoLine("Project: " + project)
	fmt.Println()

	// Show summary in a table
	clime.InfoBanner("Summary of your inputs:")

	summary := map[string]string{
		"Name":      name,
		"Email":     email,
		"Age":       fmt.Sprintf("%d", age),
		"Framework": frameworks[choice],
		"Languages": strings.Join(selectedLanguages, ", "),
		"Command":   command,
		"Project":   project,
	}

	clime.PrintKeyValueTable(summary)
	fmt.Println()

	// Show final message
	finalMessage := fmt.Sprintf("Thank you %s! All your information has been collected successfully. This demonstrates how easy it is to create beautiful interactive CLI applications with Clime!", name)

	clime.NewBox().
		WithTitle("Demo Complete").
		WithBorderColor(clime.Success).
		WithTitleColor(clime.Success).
		WithStyle(clime.BoxStyleRounded).
		AddText(finalMessage).
		Println()

	clime.SuccessLine("Interactive demo finished!")

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
