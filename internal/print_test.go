package internal

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/nalgeon/be"
	"github.com/nalgeon/howto/internal/ai"
)

func TestPrintUsage(t *testing.T) {
	out := &bytes.Buffer{}
	PrintUsage(out)
	got := out.String()
	be.True(t, strings.Contains(got, "Usage: howto [-h] [-v] [-run] [question]"))
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
	ver := NewVersion("1.2.3", "commit", "now")
	printVersion(out, ver, config, history)
	got := out.String()

	be.True(t, strings.Contains(got, bold("howto")+" 1.2.3 (now)"))
	be.True(t, strings.Contains(got, "## Config"))
	be.True(t, strings.Contains(got, "## Prompt"))
	be.True(t, strings.Contains(got, "## History"))
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
			be.Equal(t, out.String(), tt.want)
		})
	}
}

func Test_bold(t *testing.T) {
	got := bold("test")
	want := "\033[1mtest\033[0m"
	be.Equal(t, got, want)
}
