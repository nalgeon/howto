package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_howto(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := howto(out, ask, []string{"-h"}, history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), "Usage: howto [-h] [-v] [-run] [question]") {
			t.Errorf("Expected usage string, got %q", out.String())
		}
	})

	t.Run("version", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := howto(out, ask, []string{"-v"}, history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), bold("howto")+" head (now)") {
			t.Errorf("Expected version string, got %q", out.String())
		}
	})

	t.Run("run command", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", nil
		}
		history := &History{}
		err := howto(out, ask, []string{"-run"}, history)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "no command to run" {
			t.Errorf("Expected error %q, got %q", "no command to run", err.Error())
		}
	})

	t.Run("answer", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{}
		err := howto(out, ask, []string{"test"}, history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), bold("test command")) {
			t.Errorf("Expected command string, got %q", out.String())
		}
		if !strings.Contains(out.String(), "test explanation") {
			t.Errorf("Expected explanation string, got %q", out.String())
		}
	})

	t.Run("answer with follow up", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{messages: []string{"test"}}
		err := howto(out, ask, []string{"+test"}, history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), bold("test command")) {
			t.Errorf("Expected command string, got %q", out.String())
		}
		if !strings.Contains(out.String(), "test explanation") {
			t.Errorf("Expected explanation string, got %q", out.String())
		}
	})

	t.Run("answer with error", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", errors.New("test error")
		}
		history := &History{}
		err := howto(out, ask, []string{"test"}, history)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "test error" {
			t.Errorf("Expected error %q, got %q", "test error", err.Error())
		}
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
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), bold("test command")) {
			t.Errorf("Expected command string, got %q", out.String())
		}
		if !strings.Contains(out.String(), "test explanation") {
			t.Errorf("Expected explanation string, got %q", out.String())
		}
		if len(history.messages) != 2 {
			t.Errorf("Expected 2 messages in history, got %d", len(history.messages))
		}
	})

	t.Run("follow up", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "test command\ntest explanation", nil
		}
		history := &History{messages: []string{"test"}}
		err := answer(out, ask, "+test", history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), bold("test command")) {
			t.Errorf("Expected command string, got %q", out.String())
		}
		if !strings.Contains(out.String(), "test explanation") {
			t.Errorf("Expected explanation string, got %q", out.String())
		}
		if len(history.messages) != 3 {
			t.Errorf("Expected 3 messages in history, got %d", len(history.messages))
		}
	})

	t.Run("ask error", func(t *testing.T) {
		out := &bytes.Buffer{}
		ask := func(history []string) (string, error) {
			return "", errors.New("test error")
		}
		history := &History{}
		err := answer(out, ask, "test", history)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "test error" {
			t.Errorf("Expected error %q, got %q", "test error", err.Error())
		}
		if len(history.messages) != 1 {
			t.Errorf("Expected 1 message in history, got %d", len(history.messages))
		}
	})
}

func Test_answer_no_panic(t *testing.T) {
	// Test case with a nil io.Writer
	ask := func(history []string) (string, error) {
		return "test command\ntest explanation", nil
	}
	history := &History{}
	err := answer(io.Discard, ask, "test", history)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Test case with a nil ask function
	history = &History{}
	err = answer(io.Discard, nil, "test", history)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func Test_printAnswer(t *testing.T) {
	t.Run("command and explanation", func(t *testing.T) {
		out := &bytes.Buffer{}
		printAnswer(out, "test command\ntest explanation")
		if !strings.Contains(out.String(), bold("test command")) {
			t.Errorf("Expected command string, got %q", out.String())
		}
		if !strings.Contains(out.String(), "test explanation") {
			t.Errorf("Expected explanation string, got %q", out.String())
		}
	})

	t.Run("no explanation", func(t *testing.T) {
		out := &bytes.Buffer{}
		printAnswer(out, "test command")
		if !strings.Contains(out.String(), "test command") {
			t.Errorf("Expected command string, got %q", out.String())
		}
	})
}

func Test_runCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if os.Getenv("SKIP_EXEC_TEST") != "" {
			t.Skip("Skipping exec test")
		}

		out := &bytes.Buffer{}
		history := &History{messages: []string{"test", "echo test"}}
		err := runCommand(out, history)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out.String(), "test") {
			t.Errorf("Expected output string, got %q", out.String())
		}
	})

	t.Run("no command", func(t *testing.T) {
		out := &bytes.Buffer{}
		history := &History{}
		err := runCommand(out, history)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != "no command to run" {
			t.Errorf("Expected error %q, got %q", "no command to run", err.Error())
		}
	})

	t.Run("exec error", func(t *testing.T) {
		if os.Getenv("SKIP_EXEC_TEST") != "" {
			t.Skip("Skipping exec test")
		}

		out := &bytes.Buffer{}
		history := &History{messages: []string{"test", "invalid command"}}
		err := runCommand(out, history)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func Test_execCommand(t *testing.T) {
	if os.Getenv("SKIP_EXEC_TEST") != "" {
		t.Skip("Skipping exec test")
	}

	t.Run("success", func(t *testing.T) {
		out, err := execCommand("echo test")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !strings.Contains(out, "test") {
			t.Errorf("Expected output string, got %q", out)
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := execCommand("invalid command")
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("stderr", func(t *testing.T) {
		_, err := execCommand("false")
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

func Test_execCommand_no_panic(t *testing.T) {
	// Test case with an empty command
	_, err := execCommand("")
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Test case with a command that produces no output
	_, err = execCommand("true")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func Test_howto_integration(t *testing.T) {
	if os.Getenv("SKIP_EXEC_TEST") != "" {
		t.Skip("Skipping exec test")
	}

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

	// Initialize a new History instance.
	history := &History{}

	// Test case 1: Ask a question and check the output.
	out := &bytes.Buffer{}
	err := howto(out, ask, []string{"echo", "hello"}, history)
	if err != nil {
		t.Fatalf("Test case 1 failed: %v", err)
	}
	wantStr1 := bold("echo hello") + "\n\n" + "Prints hello to the console." + "\n"
	if out.String() != wantStr1 {
		t.Errorf("Test case 1 failed: expected %q, got %q", wantStr1, out.String())
	}

	// Test case 2: Run the last command and check the output.
	out.Reset()
	err = howto(out, ask, []string{"-run"}, history)
	if err != nil {
		t.Fatalf("Test case 2 failed: %v", err)
	}
	wantStr2 := bold("echo hello") + "\n\n" + "hello" + "\n"
	if out.String() != wantStr2 {
		t.Errorf("Test case 2 failed: expected %q, got %q", wantStr2, out.String())
	}

	// Test case 3: Ask a follow-up question and check the output.
	out.Reset()
	err = howto(out, ask, []string{"+echo", "world"}, history)
	if err != nil {
		t.Fatalf("Test case 3 failed: %v", err)
	}
	wantStr3 := bold("echo world") + "\n\n" + "Prints world to the console." + "\n"
	if out.String() != wantStr3 {
		t.Errorf("Test case 3 failed: expected %q, got %q", wantStr3, out.String())
	}

	// Test case 4: Verify the history.
	wantHistory := []string{"echo hello", "echo hello\n\nPrints hello to the console.", "echo world", "echo world\n\nPrints world to the console."}
	if !reflect.DeepEqual(history.messages, wantHistory) {
		t.Errorf("Test case 4 failed: expected history %v, got %v", wantHistory, history.messages)
	}
}
