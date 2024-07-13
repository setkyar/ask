package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/setkyar/ask/internal/config"
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIProvider struct {
	config       config.OpenAIConfig
	conversation []OpenAIMessage
}

func NewOpenAIProvider(cfg config.OpenAIConfig) *OpenAIProvider {
	return &OpenAIProvider{
		config: cfg,
		conversation: []OpenAIMessage{
			{Role: "system", Content: cfg.SystemMessage},
		},
	}
}

func (o *OpenAIProvider) GenerateResponse(prompt string) (string, error) {
	o.conversation = append(o.conversation, OpenAIMessage{Role: "user", Content: prompt})

	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":    o.config.Model,
		"messages": o.conversation,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message not found in response")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content not found in message")
	}

	// Add the assistant's response to the conversation
	o.conversation = append(o.conversation, OpenAIMessage{Role: "assistant", Content: content})

	return content, nil
}
