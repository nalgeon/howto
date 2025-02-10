package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ollOptions represents the options for the Ollama API.
type ollOptions struct {
	Temperature float64 `json:"temperature"`
}

// ollRequest represents the request sent to the Ollama API.
type ollRequest struct {
	Model    string     `json:"model"`
	Options  ollOptions `json:"options"`
	Stream   bool       `json:"stream"`
	Messages []message  `json:"messages"`
}

// ollAnswer represents the response from the Ollama API.
type ollAnswer struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// ollama is an AI model that uses the Ollama API.
type ollama struct {
	config Config
}

// Ask sends a question to the AI and returns the answer.
func (ai ollama) Ask(history []string) (string, error) {
	messages := buildMessages(ai.config.Prompt, history)
	req, err := ai.buildReq(messages)
	if err != nil {
		return "", err
	}

	resp, err := ai.fetchResp(req)
	if err != nil {
		return "", err
	}

	return ai.parseAnswer(resp)
}

// buildReq constructs an HTTP request from the AI configuration and messages.
func (ai ollama) buildReq(messages []message) (*http.Request, error) {
	reqBody := ollRequest{
		Model:    ai.config.Model,
		Options:  ollOptions{Temperature: ai.config.Temperature},
		Stream:   false,
		Messages: messages,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ai.config.URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// fetchResp sends the HTTP request and returns the response.
func (ai ollama) fetchResp(req *http.Request) (*http.Response, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status)
	}

	return resp, nil
}

// parseAnswer extracts the answer from the HTTP response.
func (ai ollama) parseAnswer(resp *http.Response) (string, error) {
	var ans ollAnswer
	err := json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		return "", err
	}
	content := ans.Message.Content
	return strings.TrimSpace(content), nil
}
