package ai

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/setkyar/ask/internal/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AIProvider interface {
	GenerateResponse(prompt string) (string, error)
}

func StartChat(model string) {
	cfg := config.GetConfig()
	defaultModel := cfg.DefaultModel
	if model != "" {
		defaultModel = model
	}

	if defaultModel == "" {
		fmt.Println("No default AI model set. Please run 'ask --update-config' to set up your AI providers.")
		return
	}

	var activeProvider AIProvider
	switch defaultModel {
	case "claude":
		activeProvider = NewClaudeProvider(cfg.Claude)
	case "openai":
		activeProvider = NewOpenAIProvider(cfg.OpenAI)
	default:
		fmt.Printf("Unknown AI model: %s. Please run 'ask --update-config' to set up your AI providers.\n", defaultModel)
		return
	}

	fmt.Println("Chat started. Type your messages and press Enter. Type 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		scanner.Scan()
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Chat ended. Goodbye!")
			return
		}

		response, err := activeProvider.GenerateResponse(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		c := cases.Title(language.English)
		fmt.Printf("%s 🤖: %s\n", c.String(defaultModel), response)
	}
}
