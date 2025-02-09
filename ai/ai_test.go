package ai

import (
	"reflect"
	"testing"
)

func Test_buildMessages(t *testing.T) {
	prompt := "You are a helpful assistant."

	tests := []struct {
		name    string
		history []string
		want    []message
	}{
		{
			name:    "No history",
			history: []string{},
			want: []message{
				{Role: "system", Content: prompt},
			},
		},
		{
			name:    "Single user message",
			history: []string{"Hello"},
			want: []message{
				{Role: "system", Content: prompt},
				{Role: "user", Content: "Hello"},
			},
		},
		{
			name:    "User and assistant message",
			history: []string{"Hello", "Hi there!"},
			want: []message{
				{Role: "system", Content: prompt},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
			},
		},
		{
			name:    "Multiple user and assistant messages",
			history: []string{"Hello", "Hi there!", "How are you?", "I'm fine, thank you."},
			want: []message{
				{Role: "system", Content: prompt},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
				{Role: "user", Content: "How are you?"},
				{Role: "assistant", Content: "I'm fine, thank you."},
			},
		},
		{
			name:    "Odd number of messages",
			history: []string{"Hello", "Hi there!", "How are you?"},
			want: []message{
				{Role: "system", Content: prompt},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
				{Role: "user", Content: "How are you?"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMessages(prompt, tt.history)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expected: %v, got: %v", tt.want, got)
			}
		})
	}
}
