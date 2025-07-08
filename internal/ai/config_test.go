package ai

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nalgeon/be"
)

func Test_loadConfig(t *testing.T) {
	// Save current environment variables and restore them after the test.
	oldEnv := map[string]string{}
	for _, env := range os.Environ() {
		key, value, _ := strings.Cut(env, "=")
		oldEnv[key] = value
	}
	defer func() {
		os.Clearenv()
		for k, v := range oldEnv {
			_ = os.Setenv(k, v)
		}
	}()

	tests := []struct {
		name     string
		setupEnv func()
		want     Config
		wantErr  string
	}{
		{
			name: "default config",
			setupEnv: func() {
				os.Clearenv()
			},
			want: Config{
				Vendor:      defaultVendor,
				URL:         openAIURL,
				Token:       "",
				Model:       defaultModel,
				Prompt:      "", // This will be set in the test
				Temperature: defaultTemperature,
				Timeout:     defaultTimeout,
			},
		},
		{
			name: "custom config",
			setupEnv: func() {
				_ = os.Setenv("HOWTO_AI_VENDOR", "ollama")
				_ = os.Setenv("HOWTO_AI_URL", "http://localhost:12345")
				_ = os.Setenv("HOWTO_AI_TOKEN", "test_token")
				_ = os.Setenv("HOWTO_AI_MODEL", "test_model")
				_ = os.Setenv("HOWTO_AI_PROMPT", "test_prompt")
				_ = os.Setenv("HOWTO_AI_TEMPERATURE", "0.5")
				_ = os.Setenv("HOWTO_AI_TIMEOUT", "60")
			},
			want: Config{
				Vendor:      "ollama",
				URL:         "http://localhost:12345",
				Token:       "test_token",
				Model:       "test_model",
				Prompt:      "test_prompt",
				Temperature: 0.5,
				Timeout:     60 * time.Second,
			},
		},
		{
			name: "invalid temperature",
			setupEnv: func() {
				_ = os.Setenv("HOWTO_AI_TEMPERATURE", "invalid")
			},
			want: Config{
				Vendor:      defaultVendor,
				URL:         openAIURL,
				Token:       "",
				Model:       defaultModel,
				Prompt:      "", // This will be set in the test
				Temperature: defaultTemperature,
				Timeout:     defaultTimeout,
			},
		},
		{
			name: "invalid timeout",
			setupEnv: func() {
				_ = os.Setenv("HOWTO_AI_TIMEOUT", "invalid")
			},
			want: Config{
				Vendor:      defaultVendor,
				URL:         openAIURL,
				Token:       "",
				Model:       defaultModel,
				Prompt:      "", // This will be set in the test
				Temperature: defaultTemperature,
				Timeout:     defaultTimeout,
			},
		},
		{
			name: "unknown vendor",
			setupEnv: func() {
				_ = os.Setenv("HOWTO_AI_VENDOR", "unknown")
			},
			want:    Config{},
			wantErr: "unknown AI vendor: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			tt.setupEnv()

			got, err := loadConfig()

			// Set the prompt to the default prompt for comparison, since it depends on the OS.
			if tt.want.Prompt == "" {
				tt.want.Prompt = got.Prompt
			}

			if tt.wantErr == "" {
				be.Err(t, err, nil)
			} else {
				be.Err(t, err, tt.wantErr)
			}

			be.Equal(t, got, tt.want)
		})
	}
}
