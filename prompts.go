package clime

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type InputConfig struct {
	Label       string
	Placeholder string
	Default     string
	Required    bool
	Mask        bool
	Validate    func(string) error
	Transform   func(string) string
}

type ConfirmConfig struct {
	Label   string
	Default bool
}

type SelectConfig struct {
	Label    string
	Options  []string
	Default  int
	Multiple bool
}

// Input shows a text input prompt
func Input(config InputConfig) (string, error) {
	prompt := buildInputPrompt(config)
	fmt.Print(prompt)

	var input string
	var err error

	if config.Mask {
		input, err = readPassword()
	} else {
		input, err = readLine()
	}

	if err != nil {
		return "", err
	}

	if strings.TrimSpace(input) == "" && config.Default != "" {
		input = config.Default
	}

	if config.Required && strings.TrimSpace(input) == "" {
		Error.Println("This field is required")
		return Input(config)
	}

	if config.Transform != nil {
		input = config.Transform(input)
	}

	if config.Validate != nil {
		if err := config.Validate(input); err != nil {
			Error.Printf("Validation failed: %v\n", err)
			return Input(config) // Retry
		}
	}

	return input, nil
}

// Confirm shows a yes/no confirmation prompt
func Confirm(config ConfirmConfig) (bool, error) {
	defaultText := "y/N"
	if config.Default {
		defaultText = "Y/n"
	}

	prompt := fmt.Sprintf("%s (%s): ", config.Label, defaultText)
	fmt.Print(Info.Sprint("? ") + prompt)

	input, err := readLine()
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return config.Default, nil
	}

	switch input {
	case "y", "yes", "true", "1":
		return true, nil
	case "n", "no", "false", "0":
		return false, nil
	default:
		Warning.Println("Please answer yes or no")
		return Confirm(config)
	}
}

// Select shows a single selection prompt with arrow key navigation
func Select(config SelectConfig) (int, error) {
	if len(config.Options) == 0 {
		return 0, fmt.Errorf("no options provided")
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		return selectInteractive(config)
	}

	return selectFallback(config)
}

func selectInteractive(config SelectConfig) (int, error) {
	currentSelection := config.Default
	if currentSelection >= len(config.Options) {
		currentSelection = 0
	}

	HideCursor()
	defer ShowCursor()

	displaySelectOptions(config, currentSelection)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return selectFallback(config)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		b := make([]byte, 4)
		n, err := os.Stdin.Read(b)
		if err != nil {
			return 0, err
		}

		if n == 1 {
			switch b[0] {
			case 13:
				clearSelectDisplay(len(config.Options) + 2)
				fmt.Printf("%s %s\n", Info.Sprint("?"), config.Label)
				fmt.Printf("  %s %s\n", Success.Sprint("→"), config.Options[currentSelection])
				return currentSelection, nil
				
			case 27:
				if n == 1 {
					clearSelectDisplay(len(config.Options) + 2)
					return 0, fmt.Errorf("selection cancelled")
				}
				
			case 'q', 'Q':
				clearSelectDisplay(len(config.Options) + 2)
				return 0, fmt.Errorf("selection cancelled")
			}
		} else if n >= 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65:
				if currentSelection > 0 {
					currentSelection--
				} else {
					currentSelection = len(config.Options) - 1
				}
				refreshSelectDisplay(config, currentSelection)
				
			case 66:
				if currentSelection < len(config.Options)-1 {
					currentSelection++
				} else {
					currentSelection = 0
				}
				refreshSelectDisplay(config, currentSelection)
			}
		}
	}
}

func selectFallback(config SelectConfig) (int, error) {
	fmt.Println(Info.Sprint("? ") + config.Label)

	for i, option := range config.Options {
		marker := " "
		if i == config.Default {
			marker = ">"
		}
		fmt.Printf("  %s %d) %s\n", marker, i+1, option)
	}

	fmt.Print("Select (1-" + strconv.Itoa(len(config.Options)) + "): ")

	input, err := readLine()
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return config.Default, nil
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(config.Options) {
		Error.Printf("Invalid selection. Please choose a number between 1 and %d\n", len(config.Options))
		return selectFallback(config)
	}

	return selection - 1, nil
}

func displaySelectOptions(config SelectConfig, currentSelection int) {
	fmt.Printf("%s %s\n", Info.Sprint("?"), config.Label)
	fmt.Printf("%s\n", Muted.Sprint("(↑/↓ navigate, Enter select, Esc cancel)"))
	
	for i, option := range config.Options {
		if i == currentSelection {
			fmt.Printf("  %s %s\n", Success.Sprint("→"), BoldColor.Sprint(option))
		} else {
			fmt.Printf("    %s\n", option)
		}
	}
}

func refreshSelectDisplay(config SelectConfig, currentSelection int) {
	fmt.Printf("\033[%dA", len(config.Options)+2)
	fmt.Print("\033[J")
	displaySelectOptions(config, currentSelection)
}

func clearSelectDisplay(lines int) {
	fmt.Printf("\033[%dA", lines)
	fmt.Print("\033[J")
}

// MultiSelect shows a multi-selection prompt with arrow key navigation
func MultiSelect(config SelectConfig) ([]int, error) {
	if len(config.Options) == 0 {
		return nil, fmt.Errorf("no options provided")
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		return multiSelectInteractive(config)
	}

	return multiSelectFallback(config)
}

func multiSelectInteractive(config SelectConfig) ([]int, error) {
	currentSelection := 0
	selected := make(map[int]bool)

	HideCursor()
	defer ShowCursor()

	displayMultiSelectOptions(config, currentSelection, selected)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return multiSelectFallback(config)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		b := make([]byte, 4)
		n, err := os.Stdin.Read(b)
		if err != nil {
			return nil, err
		}

		if n == 1 {
			switch b[0] {
			case 13:
				clearMultiSelectDisplay(len(config.Options) + 2)
				var result []int
				for i := 0; i < len(config.Options); i++ {
					if selected[i] {
						result = append(result, i)
					}
				}
				
				fmt.Printf("%s %s\n", Info.Sprint("?"), config.Label)
				if len(result) > 0 {
					fmt.Printf("  %s Selected %d option(s)\n", Success.Sprint("→"), len(result))
				} else {
					fmt.Printf("  %s No options selected\n", Warning.Sprint("→"))
				}
				return result, nil
				
			case 27:
				if n == 1 {
					clearMultiSelectDisplay(len(config.Options) + 2)
					return nil, fmt.Errorf("selection cancelled")
				}
				
			case 32:
				selected[currentSelection] = !selected[currentSelection]
				refreshMultiSelectDisplay(config, currentSelection, selected)
				
			case 'q', 'Q':
				clearMultiSelectDisplay(len(config.Options) + 2)
				return nil, fmt.Errorf("selection cancelled")
			}
		} else if n >= 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65:
				if currentSelection > 0 {
					currentSelection--
				} else {
					currentSelection = len(config.Options) - 1
				}
				refreshMultiSelectDisplay(config, currentSelection, selected)
				
			case 66:
				if currentSelection < len(config.Options)-1 {
					currentSelection++
				} else {
					currentSelection = 0
				}
				refreshMultiSelectDisplay(config, currentSelection, selected)
			}
		}
	}
}

func multiSelectFallback(config SelectConfig) ([]int, error) {
	selected := make(map[int]bool)

	for {
		fmt.Print("\033[2J\033[H")

		fmt.Println(Info.Sprint("? ") + config.Label + " (use space to select, enter to confirm)")

		for i, option := range config.Options {
			marker := "○"
			if selected[i] {
				marker = Success.Sprint("●")
			}
			fmt.Printf("  %s %s\n", marker, option)
		}

		fmt.Println("\nPress:")
		fmt.Println("  1-" + strconv.Itoa(len(config.Options)) + ": Toggle option")
		fmt.Println("  Enter: Confirm selection")
		fmt.Println("  q: Quit")

		input, err := readLine()
		if err != nil {
			return nil, err
		}

		input = strings.TrimSpace(input)

		if input == "" {
			var result []int
			for i := range config.Options {
				if selected[i] {
					result = append(result, i)
				}
			}
			return result, nil
		}

		if input == "q" {
			return nil, fmt.Errorf("selection cancelled")
		}

		selection, err := strconv.Atoi(input)
		if err != nil || selection < 1 || selection > len(config.Options) {
			continue
		}

		index := selection - 1
		selected[index] = !selected[index]
	}
}

func displayMultiSelectOptions(config SelectConfig, currentSelection int, selected map[int]bool) {
	fmt.Printf("%s %s\n", Info.Sprint("?"), config.Label)
	fmt.Printf("%s\n", Muted.Sprint("(↑/↓ navigate, Space select, Enter confirm, Esc cancel)"))
	
	for i, option := range config.Options {
		marker := "○"
		if selected[i] {
			marker = Success.Sprint("●")
		}
		
		if i == currentSelection {
			fmt.Printf("  %s %s %s\n", Success.Sprint("→"), marker, BoldColor.Sprint(option))
		} else {
			fmt.Printf("    %s %s\n", marker, option)
		}
	}
}

func refreshMultiSelectDisplay(config SelectConfig, currentSelection int, selected map[int]bool) {
	fmt.Printf("\033[%dA", len(config.Options)+2)
	fmt.Print("\033[J")
	displayMultiSelectOptions(config, currentSelection, selected)
}

// clearMultiSelectDisplay clears the multi-selection display
func clearMultiSelectDisplay(lines int) {
	fmt.Printf("\033[%dA", lines)
	fmt.Print("\033[J")
}

// Ask prompts for a simple text input
func Ask(label string) (string, error) {
	return Input(InputConfig{
		Label: label,
	})
}

// AskRequired prompts for a required text input
func AskRequired(label string) (string, error) {
	return Input(InputConfig{
		Label:    label,
		Required: true,
	})
}

// AskWithDefault prompts for text input with a default value
func AskWithDefault(label, defaultValue string) (string, error) {
	return Input(InputConfig{
		Label:   label,
		Default: defaultValue,
	})
}

// AskPassword prompts for a masked password input
func AskPassword(label string) (string, error) {
	return Input(InputConfig{
		Label:    label,
		Mask:     true,
		Required: true,
	})
}

// AskEmail prompts for an email with validation
func AskEmail(label string) (string, error) {
	return Input(InputConfig{
		Label:    label,
		Required: true,
		Validate: func(email string) error {
			if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
				return fmt.Errorf("invalid email format")
			}
			return nil
		},
	})
}

// AskNumber prompts for a number input
func AskNumber(label string) (int, error) {
	str, err := Input(InputConfig{
		Label:    label,
		Required: true,
		Validate: func(input string) error {
			_, err := strconv.Atoi(input)
			return err
		},
	})
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

// AskConfirm prompts for a yes/no confirmation
func AskConfirm(label string) (bool, error) {
	return Confirm(ConfirmConfig{
		Label: label,
	})
}

// AskChoice prompts for a single choice from options
func AskChoice(label string, options ...string) (int, error) {
	return Select(SelectConfig{
		Label:   label,
		Options: options,
	})
}

// AskMultiChoice prompts for multiple choices from options
func AskMultiChoice(label string, options ...string) ([]int, error) {
	return MultiSelect(SelectConfig{
		Label:    label,
		Options:  options,
		Multiple: true,
	})
}

// buildInputPrompt builds the input prompt display
func buildInputPrompt(config InputConfig) string {
	prompt := Info.Sprint("? ") + config.Label

	if config.Default != "" {
		prompt += fmt.Sprintf(" (%s)", config.Default)
	}

	if config.Placeholder != "" && config.Default == "" {
		prompt += fmt.Sprintf(" [%s]", Muted.Sprint(config.Placeholder))
	}

	if config.Required {
		prompt += Error.Sprint(" *")
	}

	prompt += ": "
	return prompt
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimRightFunc(string(line), unicode.IsSpace), nil
}

func readPassword() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return readLine()
	}

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Println()
	return string(password), nil
}

func EmailValidator(email string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("email must contain @")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("email domain must contain a dot")
	}
	return nil
}

func MinLengthValidator(min int) func(string) error {
	return func(input string) error {
		if len(input) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

func MaxLengthValidator(max int) func(string) error {
	return func(input string) error {
		if len(input) > max {
			return fmt.Errorf("must be no more than %d characters", max)
		}
		return nil
	}
}

func NumberValidator(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("must be a valid number")
	}
	return nil
}

func URLValidator(url string) error {
	url = strings.ToLower(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}
