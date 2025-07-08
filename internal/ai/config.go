package ai

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

const defaultVendor = "openai"
const openAIURL = "https://api.openai.com/v1/chat/completions"
const ollamaURL = "http://localhost:11434/api/chat"
const defaultModel = "gpt-4o"
const defaultTemperature = 0
const defaultTimeout = 30 * time.Second
const defaultPrompt = `You are a command-line assistant. You help the user solve tasks using command-line tools for the given platform (%s).

In your answer, the first line MUST be the suggested command. Do NOT use Markdown or any other formatting. Print the command in plain text WITHOUT any surrounding text.

The second line must be blank. The third line must contain a brief explanation of the command.

If you suggest multiple commands connected with pipes, you MUST provide separate explanations for each command. Print each explanation on a separate line.`

// Config describes the AI configuration.
type Config struct {
	Vendor      string
	URL         string
	Token       string
	Model       string
	Prompt      string
	Temperature float64
	Timeout     time.Duration
}

// loadConfig reads the AI configuration from environment variables.
func loadConfig() (Config, error) {
	vendor := os.Getenv("HOWTO_AI_VENDOR")
	if vendor == "" {
		vendor = defaultVendor
	}

	url := os.Getenv("HOWTO_AI_URL")
	if url == "" {
		switch vendor {
		case "openai":
			url = openAIURL
		case "ollama":
			url = ollamaURL
		default:
			err := fmt.Errorf("unknown AI vendor: %s", vendor)
			return Config{}, err
		}
	}

	token := os.Getenv("HOWTO_AI_TOKEN")

	model := os.Getenv("HOWTO_AI_MODEL")
	if model == "" {
		model = defaultModel
	}

	prompt := os.Getenv("HOWTO_AI_PROMPT")
	if prompt == "" {
		prompt = fmt.Sprintf(defaultPrompt, runtime.GOOS)
	}

	temp, err := strconv.ParseFloat(os.Getenv("HOWTO_AI_TEMPERATURE"), 64)
	if err != nil {
		temp = defaultTemperature
	}

	var timeout time.Duration
	timeoutSec, err := strconv.Atoi(os.Getenv("HOWTO_AI_TIMEOUT"))
	if err == nil {
		timeout = time.Duration(timeoutSec) * time.Second
	} else {
		timeout = defaultTimeout
	}

	return Config{
		Vendor:      vendor,
		URL:         url,
		Token:       token,
		Model:       model,
		Prompt:      prompt,
		Temperature: temp,
		Timeout:     timeout,
	}, nil
}
