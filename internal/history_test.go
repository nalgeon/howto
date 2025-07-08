package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nalgeon/be"
)

func TestHistory_Add(t *testing.T) {
	h := &History{}
	h.Add("test message")
	be.Equal(t, len(h.messages), 1)
	be.Equal(t, h.messages[0], "test message")
}

func TestHistory_Clear(t *testing.T) {
	h := &History{messages: []string{"test"}}
	h.Clear()
	be.Equal(t, len(h.messages), 0)
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
			be.Equal(t, got, tt.want)
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
			be.Equal(t, out.String(), tt.want)
		})
	}
}

func TestHistory_Save(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// LoadHistory
	hist, err := loadHistory(path)
	be.Err(t, err, nil)
	be.Equal(t, hist.path, path)

	// Add messages
	hist.Add("test question")
	hist.Add("test answer")

	// Save
	err = hist.Save()
	be.Err(t, err, nil)

	// Load again
	hist2, err := loadHistory(path)
	be.Err(t, err, nil)
	be.Equal(t, hist.messages, hist2.messages)
}

func TestHistory_Save_error(t *testing.T) {
	h := &History{path: "/invalid/path/test_history.json"}
	h.Add("test question")
	h.Add("test answer")

	err := h.Save()
	be.Err(t, err)
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
		be.Err(t, err, nil)
		wantPath := "/Library/Application Support/howto/howto-history.json"
		be.True(t, strings.HasSuffix(path, wantPath))
	})

	t.Run("linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Skipping linux test on non-linux OS")
		}
		path, err := getHistoryPath()
		be.Err(t, err, nil)
		wantPath := "/.config/howto/howto-history.json"
		be.True(t, strings.HasSuffix(path, wantPath))
	})

	t.Run("windows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping windows test on non-windows OS")
		}
		path, err := getHistoryPath()
		be.Err(t, err, nil)
		wantPath := "/howto/howto-history.json"
		be.True(t, strings.HasSuffix(path, wantPath))
	})
}

func Test_loadHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file
	err := os.WriteFile(path, []byte(`["test question", "test answer"]`), 0600)
	be.Err(t, err, nil)

	// Load history
	hist, err := loadHistory(path)
	be.Err(t, err, nil)

	// Check messages
	wantMessages := []string{"test question", "test answer"}
	be.Equal(t, hist.messages, wantMessages)
}

func Test_loadHistory_notExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Load history
	hist, err := loadHistory(path)
	be.Err(t, err, nil)

	// Check messages
	be.Equal(t, len(hist.messages), 0)
}

func Test_loadHistory_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file with invalid JSON
	err := os.WriteFile(path, []byte(`invalid json`), 0600)
	be.Err(t, err, nil)

	// Load history
	_, err = loadHistory(path)
	be.Err(t, err, "invalid character")
}

func Test_loadHistory_readonly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_history.json")

	// Create a history file
	err := os.WriteFile(path, []byte(`["test question", "test answer"]`), 0600)
	be.Err(t, err, nil)

	// Make the file read-only
	err = os.Chmod(path, 0200)
	be.Err(t, err, nil)
	defer func() { _ = os.Chmod(path, 0600) }()

	// Load history
	_, err = loadHistory(path)
	be.Err(t, err)
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
		be.Err(t, err, nil)

		err = os.MkdirAll(filepath.Dir(historyPath), 0700)
		be.Err(t, err, nil)

		err = os.WriteFile(historyPath, []byte(`["test question", "test answer"]`), 0600)
		be.Err(t, err, nil)

		// Load history.
		hist, err := LoadHistory()
		be.Err(t, err, nil)

		// Check messages.
		wantMessages := []string{"test question", "test answer"}
		be.Equal(t, hist.messages, wantMessages)
	})
}
