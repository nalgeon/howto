package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/nalgeon/howto/internal/ai"
)

// PrintUsage prints usage information.
func PrintUsage(out io.Writer) {
	fmt.Fprintln(out, "Usage: howto [-h] [-v] [-run] [question]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "A humble command-line assistant.")
	fmt.Fprintln(out, "See", underlined("https://github.com/nalgeon/howto"), "for details.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	fmt.Fprintln(out, "  -h, --help      Show this help message and exit")
	fmt.Fprintln(out, "  -v, --version   Show version information and exit")
	fmt.Fprintln(out, "  -run            Run the last suggested command")
	fmt.Fprintln(out, "  question        Describe the task to get a command suggestion")
	fmt.Fprintln(out, "                  Use '+' to ask a follow up question")
}

// printVersion prints version, configuration, and history information.
func printVersion(out io.Writer, ver Version, config ai.Config, history *History) {
	fmt.Fprintln(out, bold("howto"), ver.String())
	fmt.Fprintln(out)
	fmt.Fprintln(out, bold("## Config"))
	fmt.Fprintln(out, "- Vendor:", config.Vendor)
	fmt.Fprintln(out, "- URL:", config.URL)
	if config.Token == "" {
		fmt.Fprintln(out, "- Token: (empty)")
	} else {
		fmt.Fprintln(out, "- Token: ***")
	}
	fmt.Fprintln(out, "- Model:", config.Model)
	fmt.Fprintln(out, "- Temperature:", config.Temperature)
	fmt.Fprintln(out, "- Timeout:", config.Timeout)
	fmt.Fprintln(out)
	fmt.Fprintln(out, bold("## Prompt"))
	printWrapped(out, config.Prompt, 80)
	fmt.Fprintln(out)
	fmt.Fprintln(out, bold("## History"))
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
		fmt.Fprintln(out, text.String())
	}
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func underlined(s string) string {
	return "\033[4m" + s + "\033[0m"
}
