package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const API_ENDPOINT = "https://api.openai.com/v1/chat/completions"
const YOU = "You: "
const AI = "AI 🤖: "

var messageHistory []Message

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type completionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type completionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type Settings struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
	Role   string `json:"role"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join(homeDir, ".openai")

	settings := loadSettings(filePath)

	var question string
	var recursive, settingsOption bool
	flag.StringVar(&question, "q", "", "The question to ask the AI")
	flag.BoolVar(&recursive, "r", false, "Ask the AI a question recursively")
	flag.BoolVar(&settingsOption, "s", false, "Modify the existing api-key and model")
	flag.Parse()

	if settingsOption {
		settings = modifySettings(filePath)
	}

	if recursive {
		messageHistory = append(messageHistory, Message{Role: "system", Content: settings.Role})
		for {
			question = askUser()

			if question == "" || question == "exit" {
				fmt.Println("Bye!")
				os.Exit(0)
			}

			if question == "clear" {
				messageHistory = nil
				fmt.Print("\033[H\033[2J") // clear screen
				fmt.Println(AI, "Message history cleared.")
				continue
			}

			if question == "role" {
				fmt.Print("\033[H\033[2J") // clear screen
				messageHistory = nil
				fmt.Println(AI, "Please enter the role:")
				fmt.Scanln(YOU, &question)
				messageHistory = append(messageHistory, Message{Role: "system", Content: question})
				continue
			}

			answer := askAI(settings.APIKey, settings.Model, question)
			fmt.Println(AI, strings.TrimSpace(answer))
			fmt.Println()
		}
	}

	if question == "" {
		question = askUser()
	}

	answer := askAI(settings.APIKey, settings.Model, question)
	fmt.Println(AI, strings.TrimSpace(answer))
}

func loadSettings(filePath string) Settings {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return Settings{}
	}

	var settings Settings
	json.Unmarshal(contents, &settings)

	// check if the settings are valid
	if settings.APIKey == "" || settings.Model == "" {
		return modifySettings(filePath)
	}

	return settings
}

func modifySettings(filePath string) Settings {
	var settings Settings

	fmt.Println("Please enter your OpenAI API key:")
	// Read the API key as hidden input
	apiKeyBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Failed to read API key:", err)
	}
	settings.APIKey = string(apiKeyBytes)

	validModels := map[string]bool{
		"gpt-3.5-turbo": true,
		"gpt-4.0-turbo": true,
	}

	for {
		fmt.Println("Please choose either gpt-3.5-turbo or gpt-4.0-turbo.:")
		fmt.Scanln(&settings.Model)

		// Check if the entered model is valid
		if validModels[settings.Model] {
			break // Exit the loop if the model is valid
		} else {
			fmt.Println("Invalid model. Please choose either gpt-3.5-turbo or gpt-4.0-turbo.")
		}
	}

	if settings.Role == "" {
		fmt.Println("Please enter the role (Default is 'You are a helpful assistant'):")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			settings.Role = scanner.Text()
		}

		if settings.Role == "" {
			settings.Role = "You are a helpful assistant"
		}
	}

	contents, _ := json.Marshal(settings)
	os.WriteFile(filePath, contents, 0600)

	return settings
}

func askUser() string {
	fmt.Print(YOU)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func askAI(apiKey, model, question string) string {
	messageHistory = append(messageHistory, Message{Role: "user", Content: question})

	requestBody := chatRequest{
		Model:    model,
		Messages: messageHistory,
	}

	jsonValue, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", API_ENDPOINT, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response chatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Status code is not 200 and it's ", resp.StatusCode)
		return ""
	}

	messageHistory = append(messageHistory, Message{Role: "assistant", Content: response.Choices[0].Message.Content})

	return response.Choices[0].Message.Content
}
