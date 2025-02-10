package ai

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestOllama_Ask(t *testing.T) {
	config := Config{
		Vendor: "ollama",
		// We don't need token for olama
		Token:       "",
		Model:       "qwen2.5-coder:1.5b",
		Prompt:      "You are a test assistant.",
		Temperature: 0.7,
		Timeout:     30 * time.Second,
	}
	history := []string{"Hello", "Hi there!"}

	t.Run("successful response", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": {"content": "I'm doing great!"}}`)),
				Header:     make(http.Header),
			}
		})

		ai := ollama{config}

		answer, err := ai.Ask(history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if answer != "I'm doing great!" {
			t.Errorf("Expected answer: I'm doing great!, got: %s", answer)
		}
	})
}
