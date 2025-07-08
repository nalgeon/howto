package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nalgeon/be"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestOpenAI_Ask(t *testing.T) {
	config := Config{
		Vendor:      "openai",
		URL:         "https://test.com/v1/chat/completions",
		Token:       "test_token",
		Model:       "gpt-4",
		Prompt:      "You are a test assistant.",
		Temperature: 0.7,
		Timeout:     30 * time.Second,
	}

	history := []string{"Hello", "Hi there!"}

	t.Run("successful", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			responseBody := `{"choices": [{"message": {"content": "I'm doing great!"}}]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		})

		ai := openai{config}
		answer, err := ai.Ask(history)
		be.Err(t, err, nil)
		be.Equal(t, answer, "I'm doing great!")
	})

	t.Run("missing token", func(t *testing.T) {
		ai := openai{Config{Token: ""}}
		_, err := ai.Ask([]string{})
		be.Err(t, err, errMissingToken)
	})

	t.Run("http error", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "500 Internal Server Error",
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
		})

		ai := openai{config}
		_, err := ai.Ask(history)
		be.Err(t, err, "http status: 500 Internal Server Error")
	})

	t.Run("json decode error", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				Header:     make(http.Header),
			}
		})

		ai := openai{config}
		_, err := ai.Ask(history)
		be.Err(t, err, "invalid character")
	})

	t.Run("no answer", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			responseBody := `{"choices": []}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		})

		ai := openai{config}
		_, err := ai.Ask(history)
		be.Err(t, err, "no answer")
	})
}

func TestOpenAI_buildReq(t *testing.T) {
	config := Config{
		Vendor:      "openai",
		URL:         "https://test.com/v1/chat/completions",
		Token:       "test_token",
		Model:       "gpt-4",
		Prompt:      "You are a test assistant.",
		Temperature: 0.7,
		Timeout:     30 * time.Second,
	}
	ai := openai{config}
	messages := []message{{Role: "user", Content: "hello"}}

	req, err := ai.buildReq(messages)
	be.Err(t, err, nil)
	be.Equal(t, req.Method, http.MethodPost)
	be.Equal(t, req.URL.String(), config.URL)
	be.Equal(t, req.Header.Get("Content-Type"), "application/json")
	be.Equal(t, req.Header.Get("Authorization"), "Bearer "+config.Token)

	bodyBytes, err := io.ReadAll(req.Body)
	be.Err(t, err, nil)

	var requestBody oaiRequest
	err = json.Unmarshal(bodyBytes, &requestBody)
	be.Err(t, err, nil)

	expectedRequestBody := oaiRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
	}
	be.Equal(t, requestBody, expectedRequestBody)
}
