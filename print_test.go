package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/nalgeon/howto/ai"
)

func Test_printUsage(t *testing.T) {
	out := &bytes.Buffer{}
	printUsage(out)
	got := out.String()

	if !strings.Contains(got, "Usage: howto [-h] [-v] [-run] [question]") {
		t.Errorf("Expected usage string, got %q", got)
	}
}

func Test_printVersion(t *testing.T) {
	config := ai.Config{
		Vendor:      "test_vendor",
		URL:         "test_url",
		Token:       "test_token",
		Model:       "test_model",
		Prompt:      "test_prompt",
		Temperature: 0.5,
		Timeout:     10 * time.Second,
	}
	history := &History{messages: []string{"q1", "a1"}}

	out := &bytes.Buffer{}
	printVersion(out, "1.2.3", config, history)
	got := out.String()

	if !strings.Contains(got, bold("howto")+" 1.2.3") {
		t.Errorf("Expected version string, got %q", got)
	}
	if !strings.Contains(got, "## Config") {
		t.Errorf("Expected config header, got %q", got)
	}
	if !strings.Contains(got, "## Prompt") {
		t.Errorf("Expected prompt header, got %q", got)
	}
	if !strings.Contains(got, "## History") {
		t.Errorf("Expected history header, got %q", got)
	}
}

func Test_printWrapped(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		width int
		want  string
	}{
		{
			name:  "short string",
			s:     "hello",
			width: 10,
			want:  "hello\n",
		},
		{
			name:  "string longer than width",
			s:     "hello world",
			width: 5,
			want:  "hello\nworld\n",
		},
		{
			name:  "string with newline",
			s:     "hello\nworld",
			width: 10,
			want:  "hello\nworld\n",
		},
		{
			name:  "long word",
			s:     "thisisaverylongword",
			width: 10,
			want:  "thisisaverylongword\n",
		},
		{
			name:  "multiple lines",
			s:     "this is a long line\nand another one",
			width: 10,
			want:  "this is a\nlong line\nand\nanother\none\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printWrapped(out, tt.s, tt.width)
			if out.String() != tt.want {
				t.Errorf("Expected %q, got %q", tt.want, out.String())
			}
		})
	}
}

func Test_bold(t *testing.T) {
	got := bold("test")
	want := "\033[1mtest\033[0m"
	if got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}

func TestPrintWrapped_no_panic(t *testing.T) {
	// Test case with a nil io.Writer
	printWrapped(io.Discard, "test string", 10)

	// Test case with an empty string
	printWrapped(io.Discard, "", 10)

	// Test case with zero width
	printWrapped(io.Discard, "test string", 0)
}
