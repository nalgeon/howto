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
	ai := ollama{config}

	t.Run("successful response", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"choices": [{"message": {"content": "test"}}]}`)),
				Header:     make(http.Header),
			}
		})

		history := []string{}

		_, err := ai.Ask(history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}
