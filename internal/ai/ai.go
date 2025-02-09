// Package ai is responsible for interacting with the cloud AI
// or local Ollama models. It provides a simple vender-agnostic
// interface to ask questions and get answers.
package ai

import (
	"fmt"
	"net/http"
	"os"
)

// AskFunc is a function that sends a question to the AI.
type AskFunc func(history []string) (string, error)

// Ask sends a question to the AI and returns the answer.
// It uses the configuration prompt and conversation history
// to create a message for the AI.
// Ask is the main interface of the ai package.
var Ask AskFunc

// Conf describes the AI configuration.
var Conf Config

// HTTP client used to make requests to the AI.
var httpClient *http.Client

// message represents a single message in the conversation.
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func init() {
	// Load the configuration.
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Conf = config

	// Set the Ask function based on the vendor.
	switch config.Vendor {
	case "openai":
		Ask = openai{config}.Ask
	default:
		fmt.Println("Unknown AI vendor:", config.Vendor)
		os.Exit(1)
	}

	// Create an HTTP client with a timeout.
	httpClient = &http.Client{
		Timeout: config.Timeout,
	}
}

// buildMessages constructs a list of messages from the prompt
// and the conversation history (a sequence of user and assistant messages).
func buildMessages(prompt string, history []string) []message {
	var messages []message
	messages = append(messages, message{Role: "system", Content: prompt})
	for i := 0; i < len(history); i += 2 {
		messages = append(messages, message{Role: "user", Content: history[i]})
		if i+1 < len(history) {
			messages = append(messages, message{Role: "assistant", Content: history[i+1]})
		}
	}
	return messages
}
