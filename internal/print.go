package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/nalgeon/howto/internal/ai"
)

// PrintUsage prints usage information.
func PrintUsage(out io.Writer) {
	fprintln(out, "Usage: howto [-h] [-v] [-run] [question]")
	fprintln(out)
	fprintln(out, "A humble command-line assistant.")
	fprintln(out, "See", underlined("https://github.com/nalgeon/howto"), "for details.")
	fprintln(out)
	fprintln(out, "Options:")
	fprintln(out, "  -h, --help      Show this help message and exit")
	fprintln(out, "  -v, --version   Show version information and exit")
	fprintln(out, "  -run            Run the last suggested command")
	fprintln(out, "  question        Describe the task to get a command suggestion")
	fprintln(out, "                  Use '+' to ask a follow up question")
}

// printVersion prints version, configuration, and history information.
func printVersion(out io.Writer, ver Version, config ai.Config, history *History) {
	fprintln(out, bold("howto"), ver.String())
	fprintln(out)
	fprintln(out, bold("## Config"))
	fprintln(out, "- Vendor:", config.Vendor)
	fprintln(out, "- URL:", config.URL)
	if config.Token == "" {
		fprintln(out, "- Token: (empty)")
	} else {
		fprintln(out, "- Token: ***")
	}
	fprintln(out, "- Model:", config.Model)
	fprintln(out, "- Temperature:", config.Temperature)
	fprintln(out, "- Timeout:", config.Timeout)
	fprintln(out)
	fprintln(out, bold("## Prompt"))
	printWrapped(out, config.Prompt, 80)
	fprintln(out)
	fprintln(out, bold("## History"))
	history.Print(out)
}

// printWrapped prints a string (can be multiple lines)
// to stdout, hard-wrapping each line at the specified width.
func printWrapped(out io.Writer, s string, width int) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		var text strings.Builder
		var lineLen int

		for _, word := range words {
			wordLen := len(word)
			if lineLen == 0 {
				text.WriteString(word)
				lineLen += wordLen
			} else if lineLen+wordLen+1 <= width {
				text.WriteString(" ")
				text.WriteString(word)
				lineLen += wordLen + 1
			} else {
				text.WriteString("\n")
				text.WriteString(word)
				lineLen = wordLen
			}
		}
		fprintln(out, text.String())
	}
}

func fprintln(out io.Writer, args ...any) {
	_, _ = fmt.Fprintln(out, args...)
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func underlined(s string) string {
	return "\033[4m" + s + "\033[0m"
}
