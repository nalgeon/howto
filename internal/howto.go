package internal

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/nalgeon/howto/internal/ai"
)

// howto implements the howto command.
// Uses the given ask function to get an answer from the AI.
// Prints all output to the given writer.
func Howto(out io.Writer, ask ai.AskFunc, ver Version, args []string, history *History) error {
	input := strings.Join(args, " ")

	var err error
	switch input {
	case "-h", "--help":
		PrintUsage(out)
	case "-v", "--version":
		printVersion(out, ver, ai.Conf, history)
	case "-run":
		err = runCommand(out, history)
	default:
		err = answer(out, ask, input, history)
	}

	if err != nil {
		return err
	}

	return history.Save()
}

// answer asks the AI a question and prints the answer.
func answer(out io.Writer, ask ai.AskFunc, input string, history *History) error {
	if ask == nil {
		return fmt.Errorf("ask function is not set")
	}

	if strings.HasPrefix(input, "+") {
		input = strings.TrimSpace(input[1:])
	} else {
		history.Clear()
	}

	history.Add(input)
	answer, err := ask(history.messages)
	if err != nil {
		return err
	}

	answer = removeFences(answer)
	printAnswer(out, answer)
	history.Add(answer)
	return nil
}

// removeFences removes code fences from the answer.
func removeFences(s string) string {
	// Exclude the lines with code fences.
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s
	}
	var filtered []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "```") {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n")
}

// printAnswer prints the answer to the output.
func printAnswer(out io.Writer, answer string) {
	command, rest, ok := strings.Cut(answer, "\n")
	if !ok {
		printWrapped(out, answer, 80)
		return
	}
	_, _ = fmt.Fprintln(out, bold(command))
	printWrapped(out, rest, 80)
}

// runCommand runs the last suggested command.
func runCommand(out io.Writer, history *History) error {
	cmd := history.LastCommand()
	if cmd == "" {
		return fmt.Errorf("no command to run")
	}

	_, _ = fmt.Fprintln(out, bold(cmd))
	_, _ = fmt.Fprintln(out)
	output, err := execCommand(cmd)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, output)
	return nil
}

// execCommand executes a shell command and returns the output.
func execCommand(command string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("empty command")
	}

	var outb, errb bytes.Buffer

	// Use the shell to execute the command and avoid parsing the arguments.
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		stderr := strings.TrimSpace(errb.String())
		if stderr != "" {
			return "", fmt.Errorf("%s", stderr)
		}
		return "", err
	}

	return strings.TrimSpace(outb.String()), nil
}
