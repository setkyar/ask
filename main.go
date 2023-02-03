package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Set the API endpoint
const apiEndpoint = "https://api.openai.com/v1/completions"

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

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join(homeDir, ".openai")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}

	// File exists, read its contents
	contents, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiKey := strings.TrimSpace(string(contents))

	// Check if the API key is empty
	if apiKey == "" {
		// Ask the user to input the API key
		fmt.Println("Please enter your OpenAI API key:")
		fmt.Scanln(&apiKey)

		// Write the API key to the file
		err := os.WriteFile(filePath, []byte(apiKey), 0600)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var question string
	var recursive bool
	flag.StringVar(&question, "q", "", "The question to ask the AI")
	flag.BoolVar(&recursive, "r", false, "Ask the AI a question recursively")
	flag.Parse()

	// Check if the question is empty
	if question == "" {
		// Ask the user to input the question
		scanner := bufio.NewScanner(os.Stdin)

		// Ask the user to input the question
		fmt.Printf("You:")

		if scanner.Scan() {
			question = scanner.Text()
		}
	}

	if recursive {
		// Ask the AI a question recursively
		for {
			if question == "" || question == "exit" {
				fmt.Println("Bye!")
				break
			}

			answer := askAI(apiKey, question)
			fmt.Println("AI: ", strings.TrimSpace(answer))
			fmt.Println()

			// Ask the user to input the question
			fmt.Printf("You:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				question = strings.TrimSpace(scanner.Text())
			}
		}
	}
}

func askAI(apiKey string, question string) string {
	// Set the request body
	requestBody := completionRequest{
		Model:       "text-davinci-003",
		Prompt:      question,
		MaxTokens:   1024,
		Temperature: 0.5,
	}
	jsonValue, _ := json.Marshal(requestBody)

	// Set the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the response
	var response completionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	// check the status code
	if resp.StatusCode != 200 {
		log.Fatal("Status code is not 200 and it's ", resp.StatusCode)
		return ""
	}

	// Return the response
	return response.Choices[0].Text
}
