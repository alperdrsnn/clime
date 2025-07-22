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

func AutoComplete(config AutoCompleteConfig) (string, error) {
	if config.MaxResults == 0 {
		config.MaxResults = 10
	}
	if config.MinLength == 0 {
		config.MinLength = 1
	}

	for {
		prompt := buildAutoCompletePrompt(config)
		fmt.Print(prompt)

		input, err := readLineWithAutoComplete(config)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(input) == "" && !config.Required {
			return "", nil
		}

		if config.Required && strings.TrimSpace(input) == "" {
			Error.Println("This field is required")
			continue
		}

		if config.Transform != nil {
			input = config.Transform(input)
		}

		if config.Validate != nil {
			if err := config.Validate(input); err != nil {
				Error.Printf("Validation failed: %v\n", err)
				continue
			}
		}

		return input, nil
	}
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
	showingSuggestions := false

	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return "", err
		}

		char := b[0]

		switch char {
		case 13:
			if showingSuggestions {
				clearSuggestions(len(suggestions))
				showingSuggestions = false
			}
			fmt.Println()
			return input.String(), nil

		case 127, 8:
			if input.Len() > 0 {
				inputStr := input.String()
				input.Reset()
				input.WriteString(inputStr[:len(inputStr)-1])

				fmt.Print("\b \b")

				if showingSuggestions {
					clearSuggestions(len(suggestions))
				}
				suggestions = findSuggestions(input.String(), config)
				if len(suggestions) > 0 && input.Len() >= config.MinLength {
					displaySuggestions(suggestions)
					showingSuggestions = true
				} else {
					showingSuggestions = false
				}
			}

		case 9:
			if len(suggestions) > 0 {
				if showingSuggestions {
					clearSuggestions(len(suggestions))
				}

				suggestion := suggestions[0].Value
				input.Reset()
				input.WriteString(suggestion)

				clearCurrentLine()
				prompt := buildAutoCompletePrompt(config)
				fmt.Print(prompt + suggestion)

				showingSuggestions = false
			}

		case 27:
			nextB := make([]byte, 2)
			os.Stdin.Read(nextB)

		default:
			if char >= 32 && char <= 126 {
				input.WriteByte(char)
				fmt.Printf("%c", char)

				if showingSuggestions {
					clearSuggestions(len(suggestions))
				}
				suggestions = findSuggestions(input.String(), config)
				if len(suggestions) > 0 && input.Len() >= config.MinLength {
					displaySuggestions(suggestions)
					showingSuggestions = true
				} else {
					showingSuggestions = false
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

// displaySuggestions displays autocomplete suggestions
func displaySuggestions(suggestions []AutoCompleteResult) {
	fmt.Println()

	for i, suggestion := range suggestions {
		prefix := "  "
		if i == 0 {
			prefix = Muted.Sprint("▶ ")
		}
		fmt.Printf("%s%s\n", prefix, DimColor.Sprint(suggestion.Value))
	}

	moveCursorUp(len(suggestions) + 1)
}

// clearSuggestions clears the displayed suggestions
func clearSuggestions(count int) {
	if count == 0 {
		return
	}

	fmt.Printf("\n")
	for i := 0; i < count; i++ {
		fmt.Print("\033[2K")
		if i < count-1 {
			fmt.Print("\033[B")
		}
	}

	fmt.Printf("\033[%dA", count)
	fmt.Print("\r")
}

// clearCurrentLine clears the current line
func clearCurrentLine() {
	fmt.Print("\033[2K\r")
}

// moveCursorUp moves cursor up by n lines
func moveCursorUp(n int) {
	if n > 0 {
		fmt.Printf("\033[%dA", n)
	}
}

// buildAutoCompletePrompt builds the autocomplete prompt display
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
		Label:   label,
		Options: options,
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

// AutoCompleteBuilder provides a fluent interface for building autocomplete prompts
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
