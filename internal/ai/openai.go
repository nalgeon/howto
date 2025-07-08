package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var errMissingToken = fmt.Errorf(`set HOWTO_AI_TOKEN to your AI vendor API key.
See https://github.com/nalgeon/howto#readme for details`)

// oaiRequest represents the request sent to the OpenAI-compatible API.
type oaiRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// oaiAnswer represents the response from the OpenAI-compatible API.
type oaiAnswer struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// openai is an AI model that uses the OpenAI-compatible API.
type openai struct {
	config Config
}

// Ask sends a question to the AI and returns the answer.
func (ai openai) Ask(history []string) (string, error) {
	if ai.config.Token == "" {
		return "", errMissingToken
	}

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
func (ai openai) buildReq(messages []message) (*http.Request, error) {
	reqBody := oaiRequest{
		Model:       ai.config.Model,
		Messages:    messages,
		Temperature: ai.config.Temperature,
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
	req.Header.Set("Authorization", "Bearer "+ai.config.Token)

	return req, nil
}

// fetchResp sends the HTTP request and returns the response.
func (ai openai) fetchResp(req *http.Request) (*http.Response, error) {
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
func (ai openai) parseAnswer(resp *http.Response) (string, error) {
	var ans oaiAnswer
	err := json.NewDecoder(resp.Body).Decode(&ans)
	if err != nil {
		return "", err
	}

	if len(ans.Choices) > 0 {
		return ans.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no answer")
}
