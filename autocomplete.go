package clime

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"sort"
	"strings"
)

type AutoCompleteConfig struct {
	Label         string
	Placeholder   string
	Options       []string
	MinLength     int
	MaxResults    int
	CaseSensitive bool
	FuzzyMatch    bool
	Required      bool
	Validate      func(string) error
	Transform     func(string) string
}

type AutoCompleteResult struct {
	Value string
	Score int
	Index int
}

// AutoComplete prompts for input with autocomplete functionality
func AutoComplete(config AutoCompleteConfig) (string, error) {
	if config.MaxResults == 0 {
		config.MaxResults = 10
	}
	if config.MinLength < 0 {
		config.MinLength = 0
	}

	prompt := buildAutoCompletePrompt(config)
	fmt.Print(prompt)

	input, err := readLineWithAutoComplete(config)
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	if input == "" && !config.Required {
		return "", nil
	}

	if config.Required && input == "" {
		return "", fmt.Errorf("this field is required")
	}

	if config.Transform != nil {
		input = config.Transform(input)
	}

	if config.Validate != nil {
		if err := config.Validate(input); err != nil {
			return "", err
		}
	}

	return input, nil
}

// readLineWithAutoComplete reads input with autocomplete functionality
func readLineWithAutoComplete(config AutoCompleteConfig) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return readLine()
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return readLine()
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var input strings.Builder
	var suggestions []AutoCompleteResult
	selectedSuggestion := 0
	showingSuggestions := false

	redrawLine := func() {
		if showingSuggestions {
			clearAutoCompleteSuggestions(len(suggestions))
			showingSuggestions = false
		}
		
		suggestions = findSuggestions(input.String(), config)
		if len(suggestions) > 0 && input.Len() >= config.MinLength {
			if selectedSuggestion >= len(suggestions) {
				selectedSuggestion = 0
			}
			showSuggestions(suggestions, selectedSuggestion, input.String())
			showingSuggestions = true
		}
	}

	for {
		b := make([]byte, 4)
		n, err := os.Stdin.Read(b)
		if err != nil {
			return "", err
		}

		if n == 1 {
			switch b[0] {
			case 13:
				if showingSuggestions {
					clearAutoCompleteSuggestions(len(suggestions))
				}
				fmt.Println()
				return input.String(), nil

			case 127, 8:
				if input.Len() > 0 {
					inputStr := input.String()
					input.Reset()
					input.WriteString(inputStr[:len(inputStr)-1])
					
					fmt.Print("\b \b")
					selectedSuggestion = 0
					redrawLine()
				}

			case 9:
				if showingSuggestions && len(suggestions) > 0 {
					clearAutoCompleteSuggestions(len(suggestions))
					showingSuggestions = false
					
					backspaces := input.Len()
					input.Reset()
					input.WriteString(suggestions[selectedSuggestion].Value)
					
					for i := 0; i < backspaces; i++ {
						fmt.Print("\b")
					}
					fmt.Print(input.String())
				}

			case 27:
				continue

			default:
				if b[0] >= 32 && b[0] <= 126 {
					input.WriteByte(b[0])
					fmt.Printf("%c", b[0])
					selectedSuggestion = 0
					redrawLine()
				}
			}
		} else if n >= 3 && b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65:
				if showingSuggestions && len(suggestions) > 0 {
					if selectedSuggestion > 0 {
						selectedSuggestion--
					} else {
						selectedSuggestion = len(suggestions) - 1
					}
					clearAutoCompleteSuggestions(len(suggestions))
					showSuggestions(suggestions, selectedSuggestion, input.String())
				}
			case 66:
				if showingSuggestions && len(suggestions) > 0 {
					if selectedSuggestion < len(suggestions)-1 {
						selectedSuggestion++
					} else {
						selectedSuggestion = 0
					}
					clearAutoCompleteSuggestions(len(suggestions))
					showSuggestions(suggestions, selectedSuggestion, input.String())
				}
			}
		}
	}
}

// findSuggestions finds matching suggestions for the given input
func findSuggestions(input string, config AutoCompleteConfig) []AutoCompleteResult {
	if len(input) < config.MinLength || len(config.Options) == 0 {
		return nil
	}

	var results []AutoCompleteResult

	for i, option := range config.Options {
		score := calculateMatchScore(input, option, config)
		if score > 0 {
			results = append(results, AutoCompleteResult{
				Value: option,
				Score: score,
				Index: i,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > config.MaxResults {
		results = results[:config.MaxResults]
	}

	return results
}

// calculateMatchScore calculates how well an option matches the input
func calculateMatchScore(input, option string, config AutoCompleteConfig) int {
	if !config.CaseSensitive {
		input = strings.ToLower(input)
		option = strings.ToLower(option)
	}

	if config.FuzzyMatch {
		return fuzzyMatchScore(input, option)
	}

	if strings.HasPrefix(option, input) {
		return 1000 - len(option) + len(input)*10
	}

	if strings.Contains(option, input) {
		index := strings.Index(option, input)
		return 500 - index + len(input)*5
	}

	return 0
}

// fuzzyMatchScore calculates fuzzy match score
func fuzzyMatchScore(input, option string) int {
	if len(input) == 0 {
		return 0
	}

	score := 0
	inputIndex := 0
	consecutiveMatches := 0

	for _, char := range option {
		if inputIndex < len(input) && char == rune(input[inputIndex]) {
			score += 10 + consecutiveMatches
			consecutiveMatches++
			inputIndex++
		} else {
			consecutiveMatches = 0
		}
	}

	if inputIndex == len(input) {
		score += 100
	}

	score -= len(option) - len(input)

	return score
}

// displayAutoCompleteSuggestions displays autocomplete suggestions
func showSuggestions(suggestions []AutoCompleteResult, selected int, currentInput string) {
	fmt.Print("\n")
	
	for i, suggestion := range suggestions {
		if i == selected {
			fmt.Printf("  %s %s\n", Success.Sprint("→"), BoldColor.Sprint(suggestion.Value))
		} else {
			fmt.Printf("    %s\n", DimColor.Sprint(suggestion.Value))
		}
	}
	
	fmt.Printf("\033[%dA", len(suggestions)+1)
	fmt.Print("\033[999C")
}

// clearAutoCompleteSuggestions clears autocomplete suggestions
func clearAutoCompleteSuggestions(lines int) {
	if lines <= 0 {
		return
	}
	
	fmt.Print("\n")
	for i := 0; i < lines; i++ {
		fmt.Print("\033[2K")
		if i < lines-1 {
			fmt.Print("\033[B")
		}
	}
	fmt.Printf("\033[%dA", lines+1)
	fmt.Print("\033[999C")
}

// buildAutoCompletePrompt builds the autocomplete prompt
func buildAutoCompletePrompt(config AutoCompleteConfig) string {
	prompt := Info.Sprint("? ") + config.Label

	if config.Placeholder != "" {
		prompt += fmt.Sprintf(" [%s]", Muted.Sprint(config.Placeholder))
	}

	if config.Required {
		prompt += Error.Sprint(" *")
	}

	prompt += ": "
	return prompt
}

// AskWithOptions prompts for input with predefined options
func AskWithOptions(label string, options []string) (string, error) {
	return AutoComplete(AutoCompleteConfig{
		Label:         label,
		Options:       options,
		MinLength:     0,
		MaxResults:    8,
		CaseSensitive: false,
		FuzzyMatch:    true,
	})
}

// AskWithFileCompletion prompts for a file path with file completion
func AskWithFileCompletion(label string) (string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return Ask(label)
	}

	var options []string
	for _, file := range files {
		options = append(options, file.Name())
	}

	return AutoComplete(AutoCompleteConfig{
		Label:      label,
		Options:    options,
		FuzzyMatch: true,
	})
}

// AskWithCommandCompletion prompts with common command completion
func AskWithCommandCompletion(label string) (string, error) {
	commands := []string{
		"help", "version", "init", "start", "stop", "restart",
		"status", "config", "install", "uninstall", "update",
		"list", "show", "create", "delete", "edit", "copy",
		"move", "rename", "search", "find", "replace",
	}

	return AutoComplete(AutoCompleteConfig{
		Label:      label,
		Options:    commands,
		FuzzyMatch: true,
	})
}

type AutoCompleteBuilder struct {
	config AutoCompleteConfig
}

// NewAutoCompleteBuilder creates a new autocomplete builder
func NewAutoCompleteBuilder(label string) *AutoCompleteBuilder {
	return &AutoCompleteBuilder{
		config: AutoCompleteConfig{
			Label:      label,
			MaxResults: 10,
			MinLength:  1,
		},
	}
}

// WithOptions sets the autocomplete options
func (b *AutoCompleteBuilder) WithOptions(options []string) *AutoCompleteBuilder {
	b.config.Options = options
	return b
}

// WithPlaceholder sets the placeholder text
func (b *AutoCompleteBuilder) WithPlaceholder(placeholder string) *AutoCompleteBuilder {
	b.config.Placeholder = placeholder
	return b
}

// WithMinLength sets the minimum input length before showing suggestions
func (b *AutoCompleteBuilder) WithMinLength(length int) *AutoCompleteBuilder {
	b.config.MinLength = length
	return b
}

// WithMaxResults sets the maximum number of suggestions to show
func (b *AutoCompleteBuilder) WithMaxResults(max int) *AutoCompleteBuilder {
	b.config.MaxResults = max
	return b
}

// CaseSensitive enables case-sensitive matching
func (b *AutoCompleteBuilder) CaseSensitive(enabled bool) *AutoCompleteBuilder {
	b.config.CaseSensitive = enabled
	return b
}

// FuzzyMatch enables fuzzy matching
func (b *AutoCompleteBuilder) FuzzyMatch(enabled bool) *AutoCompleteBuilder {
	b.config.FuzzyMatch = enabled
	return b
}

// Required makes the input required
func (b *AutoCompleteBuilder) Required(required bool) *AutoCompleteBuilder {
	b.config.Required = required
	return b
}

// WithValidator sets a validation function
func (b *AutoCompleteBuilder) WithValidator(validator func(string) error) *AutoCompleteBuilder {
	b.config.Validate = validator
	return b
}

// WithTransformer sets a transformation function
func (b *AutoCompleteBuilder) WithTransformer(transformer func(string) string) *AutoCompleteBuilder {
	b.config.Transform = transformer
	return b
}

// Ask executes the autocomplete prompt
func (b *AutoCompleteBuilder) Ask() (string, error) {
	return AutoComplete(b.config)
}

var BooleanOptions = []string{"yes", "no", "true", "false", "y", "n"}

var ColorOptions = []string{
	"red", "green", "blue", "yellow", "cyan", "magenta", "white", "black",
	"gray", "orange", "pink", "purple", "brown", "lime", "navy", "teal",
}

var SizeOptions = []string{"small", "medium", "large", "xl", "xs", "xxl", "tiny", "huge"}

var PriorityOptions = []string{"low", "medium", "high", "critical", "urgent", "normal"}

var StatusOptions = []string{
	"active", "inactive", "pending", "completed", "failed", "cancelled",
	"draft", "published", "archived", "deleted",
}
