package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/setkyar/ask/internal/ai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	updateConfig bool
	provider     string
)

var rootCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask is a CLI application for interacting with AI providers",
	Long: `Ask is a CLI application that allows you to interact with AI providers like Claude and OpenAI.
It provides various commands to set up and use these AI services.`,
	Run: func(cmd *cobra.Command, args []string) {
		if updateConfig {
			setupAIProviders()
		} else if !viper.GetBool("setup_complete") {
			setupAIProviders()
		} else if provider != "" {
			ai.StartChat(provider)
		} else {
			ai.StartChat("")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolVar(&updateConfig, "update-config", false, "Update the configuration settings")
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "Specify AI provider (claude or openai)")
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}
	viper.SetConfigName(".ask_ai_settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
	}
}

func setupAIProviders() {
	if viper.GetBool("setup_complete") && !updateConfig {
		fmt.Println("AI providers are already set up. Use --update-config to modify settings.")
		return
	}

	fmt.Println("Welcome to Ask! Let's set up your AI providers.")

	setupClaude := promptYesNo("Do you want to set up or update Claude?")
	if setupClaude {
		claudeAPIKey := promptWithDefault("Enter your Claude API key:", viper.GetString("claude.api_key"))
		viper.Set("claude.api_key", claudeAPIKey)
		claudeVersion := promptWithDefault("Enter your Claude API Version:", viper.GetString("claude.api_version"))
		viper.Set("claude.api_version", claudeVersion)
		claudeModel := promptWithDefault("Enter your Claude Model:", viper.GetString("claude.model"))
		viper.Set("claude.model", claudeModel)
		claudeMaxToken := promptIntWithDefault("Enter your Claude max token:", viper.GetInt("claude.max_token"))
		viper.Set("claude.max_token", claudeMaxToken)
		claudeSystemMsg := promptWithDefault("Enter system message for Claude:", viper.GetString("claude.system_message"))
		viper.Set("claude.system_message", claudeSystemMsg)
	}

	setupOpenAI := promptYesNo("Do you want to set up or update OpenAI?")
	if setupOpenAI {
		openAIAPIKey := promptWithDefault("Enter your OpenAI API key:", viper.GetString("openai.api_key"))
		viper.Set("openai.api_key", openAIAPIKey)
		openAIModel := promptWithDefault("Enter your OpenAI Model:", viper.GetString("openai.model"))
		viper.Set("openai.model", openAIModel)
		openAISystemMsg := promptWithDefault("Enter system message for OpenAI:", viper.GetString("openai.system_message"))
		viper.Set("openai.system_message", openAISystemMsg)
	}
	defaultModel := viper.GetString("default_model")
	if (setupClaude && setupOpenAI) || defaultModel == "" {
		defaultModel = promptForDefaultModel()
		viper.Set("default_model", defaultModel)
	} else if setupClaude {
		viper.Set("default_model", "claude")
	} else if setupOpenAI {
		viper.Set("default_model", "openai")
	}

	viper.Set("setup_complete", true)

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".ask_ai_settings.yaml")
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}
	fmt.Println("Setup complete! Your settings have been saved.")
}

func promptWithDefault(question, defaultValue string) string {
	fmt.Printf("%s [%s]: ", question, defaultValue)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(scanner.Text())
	if answer == "" {
		return defaultValue
	}
	return answer
}

func promptIntWithDefault(question string, defaultValue int) int {
	for {
		fmt.Printf("%s [%d]: ", question, defaultValue)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := strings.TrimSpace(scanner.Text())
		if answer == "" {
			return defaultValue
		}
		if value, err := strconv.Atoi(answer); err == nil {
			return value
		}
		fmt.Println("Please enter a valid integer.")
	}
}

func promptYesNo(question string) bool {
	for {
		var answer string
		fmt.Print(question + " (y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			answer = strings.TrimSpace(scanner.Text())
		} else {
			fmt.Println("Error reading input:", scanner.Err())
			continue
		}
		switch strings.ToLower(answer) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Please answer with 'y' or 'n'.")
		}
	}
}

func promptForDefaultModel() string {
	for {
		fmt.Print("Which AI model would you like to use as default? (claude/openai): ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch answer {
		case "claude", "openai":
			return answer
		default:
			fmt.Println("Please enter either 'claude' or 'openai'.")
		}
	}
}
