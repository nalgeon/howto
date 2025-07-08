package internal

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/nalgeon/be"
)

func TestHowto(t *testing.T) {
	ver := NewVersion("1.2.3", "commit", "now")

	t.Run("help", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := Howto(out, ask, ver, []string{"-h"}, history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), "Usage: howto [-h] [-v] [-run] [question]"))
	})

	t.Run("version", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := Howto(out, ask, ver, []string{"-v"}, history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), bold("howto")+" 1.2.3 (now)"))
	})

	t.Run("run command", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := Howto(out, ask, ver, []string{"-run"}, history)
		be.Err(t, err, "no command to run")
	})

	t.Run("answer", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{}
		err := Howto(out, ask, ver, []string{"test"}, history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), bold("test command")))
		be.True(t, strings.Contains(out.String(), "test explanation"))
	})

	t.Run("answer with follow up", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{messages: []string{"test"}}
		err := Howto(out, ask, ver, []string{"+test"}, history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), bold("test command")))
		be.True(t, strings.Contains(out.String(), "test explanation"))
	})

	t.Run("answer with error", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", errors.New("test error")
		}
		history := &History{}
		err := Howto(out, ask, ver, []string{"test"}, history)
		be.Err(t, err, "test error")
	})
}

func Test_answer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{}
		err := answer(out, ask, "test", history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), bold("test command")))
		be.True(t, strings.Contains(out.String(), "test explanation"))
		be.Equal(t, len(history.messages), 2)
	})

	t.Run("follow up", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{messages: []string{"test"}}
		err := answer(out, ask, "+test", history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), bold("test command")))
		be.True(t, strings.Contains(out.String(), "test explanation"))
		be.Equal(t, len(history.messages), 3)
	})

	t.Run("ask error", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", errors.New("test error")
		}
		history := &History{}
		err := answer(out, ask, "test", history)
		be.Err(t, err, "test error")
		be.Equal(t, len(history.messages), 1)
	})
}

func Test_removeFences(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "no fences",
			s:    "test",
			want: "test",
		},
		{
			name: "single fence",
			s:    "```test```",
			want: "```test```",
		},
		{
			name: "fences with text",
			s:    "```bash\ntest\n```\nafter",
			want: "test\nafter",
		},
		{
			name: "multiple fences",
			s:    "before\n```\ntest\n```\nmid\n```\ntest\n```\nafter",
			want: "before\ntest\nmid\ntest\nafter",
		},
		{
			name: "fences with spaces",
			s:    "  ```bash\n  test\n  ```\nafter",
			want: "test\nafter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeFences(tt.s)
			be.Equal(t, got, tt.want)
		})
	}
}

func Test_runCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		out := &bytes.Buffer{}
		history := &History{messages: []string{"test", "echo test"}}
		err := runCommand(out, history)
		be.Err(t, err, nil)
		be.True(t, strings.Contains(out.String(), "test"))
	})

	t.Run("no command", func(t *testing.T) {
		out := &bytes.Buffer{}
		history := &History{}
		err := runCommand(out, history)
		be.Err(t, err, "no command to run")
	})

	t.Run("exec error", func(t *testing.T) {
		out := &bytes.Buffer{}
		history := &History{messages: []string{"test", "invalid command"}}
		err := runCommand(out, history)
		be.Err(t, err)
	})
}

func TestHowto_integration(t *testing.T) {
	// Define a mock AI ask function for testing purposes.
	ask := func(history []string) (string, error) {
		question := history[len(history)-1]
		switch question {
		case "echo hello":
			return "echo hello\n\nPrints hello to the console.", nil
		case "echo world":
			return "echo world\n\nPrints world to the console.", nil
		default:
			return "", fmt.Errorf("unexpected question: %s", question)
		}
	}

	ver := NewVersion("1.2.3", "commit", "now")
	history := &History{}

	// Test case 1: Ask a question and check the output.
	out := &bytes.Buffer{}
	err := Howto(out, ask, ver, []string{"echo", "hello"}, history)
	be.Err(t, err, nil)
	wantStr1 := bold("echo hello") + "\n\n" + "Prints hello to the console." + "\n"
	be.Equal(t, out.String(), wantStr1)

	// Test case 2: Run the last command and check the output.
	out.Reset()
	err = Howto(out, ask, ver, []string{"-run"}, history)
	be.Err(t, err, nil)
	wantStr2 := bold("echo hello") + "\n\n" + "hello" + "\n"
	be.Equal(t, out.String(), wantStr2)

	// Test case 3: Ask a follow-up question and check the output.
	out.Reset()
	err = Howto(out, ask, ver, []string{"+echo", "world"}, history)
	be.Err(t, err, nil)
	wantStr3 := bold("echo world") + "\n\n" + "Prints world to the console." + "\n"
	be.Equal(t, out.String(), wantStr3)

	// Test case 4: Verify the history.
	wantHistory := []string{"echo hello", "echo hello\n\nPrints hello to the console.", "echo world", "echo world\n\nPrints world to the console."}
	be.Equal(t, history.messages, wantHistory)
}
