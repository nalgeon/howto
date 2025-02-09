package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestHistory_Add(t *testing.T) {
	h := &History{}
	h.Add("test message")

	if len(h.messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(h.messages))
	}

	if h.messages[0] != "test message" {
		t.Errorf("Expected message %q, got %q", "test message", h.messages[0])
	}
}

func TestHistory_Clear(t *testing.T) {
	h := &History{messages: []string{"test"}}
	h.Clear()

	if len(h.messages) != 0 {
		t.Fatalf("Expected 0 messages, got %d", len(h.messages))
	}
}

func TestHistory_LastCommand(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		want     string
	}{
		{
			name:     "empty history",
			messages: []string{},
			want:     "",
		},
		{
			name:     "single message",
			messages: []string{"command\nexplanation"},
			want:     "command",
		},
		{
			name:     "multiple messages",
			messages: []string{"q1", "a1\ncmd1\nexp1", "q2", "a2\ncmd2\nexp2"},
			want:     "a2",
		},
		{
			name:     "no newline",
			messages: []string{"command"},
			want:     "command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &History{messages: tt.messages}
			got := h.LastCommand()
			if got != tt.want {
				t.Errorf("Expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestHistory_Print(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		want     string
	}{
		{
			name:     "empty history",
			messages: []string{},
			want:     "(empty)\n",
		},
		{
			name:     "single user message",
			messages: []string{"hello"},
			want:     "🧑 hello\n",
		},
		{
			name:     "single assistant message",
			messages: []string{"hello", "hi"},
			want:     "🧑 hello\n🤖 hi\n",
		},
		{
			name:     "long message",
			messages: []string{"this is a very very very very very very very long message that should be truncated"},
			want:     "🧑 this is a very very very very very very very long message that should be...\n",
		},
		{
			name:     "message with newline",
			messages: []string{"hello\nworld"},
			want:     "🧑 hello world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &History{messages: tt.messages}
			out := &bytes.Buffer{}
			h.Print(out)
			if out.String() != tt.want {
				t.Errorf("Expected %q, got %q", tt.want, out.String())
			}
		})
	}
}

func TestHistory_Save(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// LoadHistory
	hist, err := loadHistory(path)
	if err != nil {
		t.Fatalf("loadHistory error: %v", err)
	}
	if hist.path != path {
		t.Errorf("Expected path %q, got %q", path, hist.path)
	}

	// Add messages
	hist.Add("test question")
	hist.Add("test answer")

	// Save
	err = hist.Save()
	if err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// Load again
	hist2, err := loadHistory(path)
	if err != nil {
		t.Fatalf("loadHistory error: %v", err)
	}

	// Check messages
	if !reflect.DeepEqual(hist.messages, hist2.messages) {
		t.Errorf("Expected messages %v, got %v", hist.messages, hist2.messages)
	}
}

func TestHistory_Save_error(t *testing.T) {
	h := &History{path: "/invalid/path/test_history.json"}
	h.Add("test question")
	h.Add("test answer")

	err := h.Save()
	if err == nil {
		t.Fatalf("Save should return an error")
	}
}

func Test_getHistoryPath(t *testing.T) {
	// Save current environment variables and restore them after the test.
	oldEnv := map[string]string{}
	for _, env := range os.Environ() {
		key, value, _ := strings.Cut(env, "=")
		oldEnv[key] = value
	}
	defer func() {
		os.Clearenv()
		for k, v := range oldEnv {
			os.Setenv(k, v)
		}
	}()

	t.Run("darwin", func(t *testing.T) {
		if runtime.GOOS != "darwin" {
			t.Skip("Skipping darwin test on non-darwin OS")
		}
		path, err := getHistoryPath()
		if err != nil {
			t.Fatalf("getHistoryPath error: %v", err)
		}
		wantPath := "/Library/Application Support/howto/howto-history.json"
		if !strings.HasSuffix(path, wantPath) {
			t.Errorf("Expected path ...%q, got %q", wantPath, path)
		}
	})

	t.Run("linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Skipping linux test on non-linux OS")
		}
		path, err := getHistoryPath()
		if err != nil {
			t.Fatalf("getHistoryPath error: %v", err)
		}
		wantPath := "/.config/howto/howto-history.json"
		if !strings.HasSuffix(path, wantPath) {
			t.Errorf("Expected path ...%q, got %q", wantPath, path)
		}
	})

	t.Run("windows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping windows test on non-windows OS")
		}
		path, err := getHistoryPath()
		if err != nil {
			t.Fatalf("getHistoryPath error: %v", err)
		}
		wantPath := "/howto/howto-history.json"
		if !strings.HasSuffix(path, wantPath) {
			t.Errorf("Expected path ...%q, got %q", wantPath, path)
		}
	})
}

func Test_loadHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file
	err := os.WriteFile(path, []byte(`["test question", "test answer"]`), 0600)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	// Load history
	hist, err := loadHistory(path)
	if err != nil {
		t.Fatalf("loadHistory error: %v", err)
	}

	// Check messages
	wantMessages := []string{"test question", "test answer"}
	if !reflect.DeepEqual(hist.messages, wantMessages) {
		t.Errorf("Expected messages %v, got %v", wantMessages, hist.messages)
	}
}

func Test_loadHistory_notExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Load history
	hist, err := loadHistory(path)
	if err != nil {
		t.Fatalf("loadHistory error: %v", err)
	}

	// Check messages
	if len(hist.messages) != 0 {
		t.Errorf("Expected empty messages, got %v", hist.messages)
	}
}

func Test_loadHistory_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file with invalid JSON
	err := os.WriteFile(path, []byte(`invalid json`), 0600)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	// Load history
	_, err = loadHistory(path)
	if err == nil {
		t.Fatalf("loadHistory should return an error")
	}
	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("Expected error containing %q, got %q", "invalid character", err.Error())
	}
}

func Test_loadHistory_readonly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file
	err := os.WriteFile(path, []byte(`["test question", "test answer"]`), 0600)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	// Make the file read-only
	err = os.Chmod(path, 0200)
	if err != nil {
		t.Fatalf("Chmod error: %v", err)
	}
	defer func() { _ = os.Chmod(path, 0600) }()

	// Load history
	_, err = loadHistory(path)
	if err == nil {
		t.Fatalf("loadHistory should return an error")
	}
}

func TestLoadHistory(t *testing.T) {
	// Save current environment variables and restore them after the test.
	oldEnv := map[string]string{}
	for _, env := range os.Environ() {
		key, value, _ := strings.Cut(env, "=")
		oldEnv[key] = value
	}
	defer func() {
		os.Clearenv()
		for k, v := range oldEnv {
			os.Setenv(k, v)
		}
	}()

	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		os.Setenv("HOME", dir)

		// Create a dummy history file.
		historyPath, err := getHistoryPath()
		if err != nil {
			t.Fatalf("getHistoryPath error: %v", err)
		}

		err = os.MkdirAll(filepath.Dir(historyPath), 0700)
		if err != nil {
			t.Fatalf("MkdirAll error: %v", err)
		}

		err = os.WriteFile(historyPath, []byte(`["test question", "test answer"]`), 0600)
		if err != nil {
			t.Fatalf("WriteFile error: %v", err)
		}

		// Load history.
		hist, err := LoadHistory()
		if err != nil {
			t.Fatalf("LoadHistory error: %v", err)
		}

		// Check messages.
		wantMessages := []string{"test question", "test answer"}
		if !reflect.DeepEqual(hist.messages, wantMessages) {
			t.Errorf("Expected messages %v, got %v", wantMessages, hist.messages)
		}
	})
}
