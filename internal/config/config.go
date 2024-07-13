package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DefaultModel string
	OpenAI       OpenAIConfig
	Claude       ClaudeConfig
}

type OpenAIConfig struct {
	APIKey        string
	Model         string
	SystemMessage string
}

type ClaudeConfig struct {
	APIKey        string
	APIVersion    string
	Model         string
	MaxTokens     int
	SystemMessage string
}

func GetConfig() *Config {
	return &Config{
		DefaultModel: viper.GetString("default_model"),
		OpenAI: OpenAIConfig{
			APIKey:        viper.GetString("openai.api_key"),
			Model:         viper.GetString("openai.model"),
			SystemMessage: viper.GetString("openai.system_message"),
		},
		Claude: ClaudeConfig{
			APIKey:        viper.GetString("claude.api_key"),
			APIVersion:    viper.GetString("claude.api_version"),
			Model:         viper.GetString("claude.model"),
			MaxTokens:     viper.GetInt("claude.max_token"),
			SystemMessage: viper.GetString("claude.system_message"),
		},
	}
}
