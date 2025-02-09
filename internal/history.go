package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Name of the file containing the history.
const fileName = "howto-history.json"

// History represents the conversation history
// between the user and the assistant.
type History struct {
	path     string
	messages []string
}

// LoadHistory loads the conversation history from the file system.
func LoadHistory() (*History, error) {
	path, err := getHistoryPath()
	if err != nil {
		return nil, fmt.Errorf("load history: %w", err)
	}
	hist, err := loadHistory(path)
	if err != nil {
		return nil, fmt.Errorf("load history: %w", err)
	}
	return hist, nil
}

// Save saves the conversation history to the file system.
func (h *History) Save() error {
	if h.path == "" {
		// Transient history, no need to save.
		return nil
	}
	data, err := json.Marshal(h.messages)
	if err != nil {
		return fmt.Errorf("save history: %w", err)
	}
	err = os.WriteFile(h.path, data, 0600)
	if err != nil {
		return fmt.Errorf("save history: %w", err)
	}
	return nil
}

// Add adds a message to the conversation history.
func (h *History) Add(message string) {
	h.messages = append(h.messages, message)
}

// Clear clears the conversation history.
func (h *History) Clear() {
	h.messages = []string{}
}

// LastCommand returns the last command from the conversation history.
// By design, the last command is always the first line of the last message
// (which is an answer from the assistant).
func (h *History) LastCommand() string {
	if len(h.messages) == 0 {
		return ""
	}
	lastMessage := h.messages[len(h.messages)-1]
	return strings.Split(lastMessage, "\n")[0]
}

// Print prints the conversation history to stdout.
func (h *History) Print(out io.Writer) {
	if len(h.messages) == 0 {
		fmt.Fprintln(out, "(empty)")
		return
	}
	for i, message := range h.messages {
		var prefix string
		if i%2 == 0 {
			prefix = "🧑 "
		} else {
			prefix = "🤖 "
		}
		message = prefix + strings.ReplaceAll(message, "\n", " ")
		if len(message) > 80 {
			message = message[:77] + "..."
		}
		fmt.Fprintln(out, message)
	}
}

// getHistoryPath returns the path to the history file.
// Uses the OS-specific configuration directory
// with a fallback to the home directory.
func getHistoryPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "howto")
	case "linux":
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "howto")
	case "windows":
		configDir = filepath.Join(os.Getenv("AppData"), "howto")
	default:
		// Fallback to home directory if OS is not recognized
		usr, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(usr, ".howto")
	}

	// Create the config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(configDir, fileName), nil
}

// loadHistory loads the conversation history from the specified file.
func loadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{path: path}, nil
	}
	if err != nil {
		return nil, err
	}

	var messages []string
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	hist := &History{path: path, messages: messages}
	return hist, nil
}
