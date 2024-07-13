package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/setkyar/ask/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeProvider struct {
	config       config.ClaudeConfig
	conversation []Message
}

func NewClaudeProvider(cfg config.ClaudeConfig) *ClaudeProvider {
	return &ClaudeProvider{
		config:       cfg,
		conversation: []Message{},
	}
}

func (c *ClaudeProvider) GenerateResponse(prompt string) (string, error) {
	c.conversation = append(c.conversation, Message{Role: "user", Content: prompt})

	url := "https://api.anthropic.com/v1/messages"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":      c.config.Model,
		"max_tokens": c.config.MaxTokens,
		"messages":   c.conversation,
		"system":     c.config.SystemMessage,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.APIVersion)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	firstContent, ok := content[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected content format")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("text not found in response")
	}

	// Add the assistant's response to the conversation
	c.conversation = append(c.conversation, Message{Role: "assistant", Content: text})

	return text, nil
}
