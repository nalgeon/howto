package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
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

	t.Run("successful request", func(t *testing.T) {
		want := oaiRequest{
			Model:       config.Model,
			Messages:    buildMessages(config.Prompt, history),
			Temperature: config.Temperature,
		}

		wantBytes, _ := json.Marshal(want)
		wantStr := string(wantBytes)

		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.String() != config.URL {
				t.Errorf("Expected URL: %s, got: %s", config.URL, req.URL.String())
			}

			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type: application/json, got: %s", req.Header.Get("Content-Type"))
			}

			if req.Header.Get("Authorization") != "Bearer "+config.Token {
				t.Errorf("Expected Authorization: Bearer %s, got: %s", config.Token, req.Header.Get("Authorization"))
			}

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("Error reading request body: %v", err)
			}
			bodyString := string(bodyBytes)

			if bodyString != wantStr {
				t.Errorf("Expected request body: %s, got: %s", wantStr, bodyString)
			}

			responseBody := `{"choices": [{"message": {"content": "I'm doing great!"}}]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		})

		ai := openai{config}
		answer, err := ai.Ask(history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if answer != "I'm doing great!" {
			t.Errorf("Expected answer: I'm doing great!, got: %s", answer)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		ai := openai{Config{Token: ""}}
		_, err := ai.Ask([]string{})
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != errMissingToken.Error() {
			t.Fatalf("Expected error %q, got %q", errMissingToken, err)
		}
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
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "http status: 500 Internal Server Error") {
			t.Fatalf("Expected error containing %q, got %q", "http status: 500 Internal Server Error", err)
		}
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
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid character") {
			t.Fatalf("Expected error containing %q, got %q", "invalid character", err)
		}
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
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "no answer" {
			t.Fatalf("Expected error %q, got %q", "no answer", err)
		}
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
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("Expected method %q, got %q", http.MethodPost, req.Method)
	}

	if req.URL.String() != config.URL {
		t.Errorf("Expected URL %q, got %q", config.URL, req.URL.String())
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", req.Header.Get("Content-Type"))
	}

	if req.Header.Get("Authorization") != "Bearer "+config.Token {
		t.Errorf("Expected Authorization %q, got %q", "Bearer "+config.Token, req.Header.Get("Authorization"))
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}

	var requestBody oaiRequest
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		t.Fatalf("Error unmarshaling request body: %v", err)
	}

	expectedRequestBody := oaiRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
	}

	if !reflect.DeepEqual(requestBody, expectedRequestBody) {
		t.Errorf("Expected request body %v, got %v", expectedRequestBody, requestBody)
	}
}

func TestOpenAI_fetchResp(t *testing.T) {
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

	t.Run("successful response", func(t *testing.T) {
		httpClient = NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"choices": [{"message": {"content": "test"}}]}`)),
				Header:     make(http.Header),
			}
		})

		req, _ := http.NewRequest("POST", config.URL, nil)
		resp, err := ai.fetchResp(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
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

		req, _ := http.NewRequest("POST", config.URL, nil)
		_, err := ai.fetchResp(req)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "http status: 500 Internal Server Error") {
			t.Fatalf("Expected error containing %q, got %q", "http status", err.Error())
		}
	})
}

func TestOpenAI_parseAnswer(t *testing.T) {
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

	t.Run("successful parse", func(t *testing.T) {
		responseBody := `{"choices": [{"message": {"content": "test answer"}}]}`
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			Header:     make(http.Header),
		}

		answer, err := ai.parseAnswer(resp)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if answer != "test answer" {
			t.Errorf("Expected answer %q, got %q", "test answer", answer)
		}
	})

	t.Run("no answer", func(t *testing.T) {
		responseBody := `{"choices": []}`
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			Header:     make(http.Header),
		}

		_, err := ai.parseAnswer(resp)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "no answer" {
			t.Errorf("Expected error %q, got %q", "no answer", err.Error())
		}
	})

	t.Run("decode error", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
			Header:     make(http.Header),
		}

		_, err := ai.parseAnswer(resp)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("Expected error containing %q, got %q", "invalid character", err.Error())
		}
	})
}
