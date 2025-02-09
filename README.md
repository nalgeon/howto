## Howto - a humble command-line assistant

Howto helps you solve command-line tasks with AI. Describe the task, and `howto` will suggest a solution:

```text
$ howto curl example.org but print only the headers
curl -I example.org

The `curl` command is used to transfer data from or to a server.
The `-I` option tells `curl` to fetch the HTTP headers only, without the body
content.
```

Howto works with any OpenAI-compatible provider and local Ollama models (coming soon). It's a simple tool that doesn't interfere with your terminal. Not an "intelligent terminal" or anything. You ask, and howto answers. That's the deal.

```text
Usage: howto [-h] [-v] [-run] [question]

A humble command-line assistant.

Options:
  -h, --help      Show this help message and exit
  -v, --version   Show version information and exit
  -run            Run the last suggested command
  question        Describe the task to get a command suggestion
                  Use '+' to ask a follow up question
```

There are some additional features you may find useful. See the Usage section for details.

## Installation

### Brew

This method is preferred if you use Homebrew:

```text
brew tap nalgeon/howto https://github.com/nalgeon/howto
brew install howto
```

### Go install

This method is preferred if you have Go installed:

```text
go install github.com/nalgeon/howto@latest
```

### Manual

Howto is a binary executable file (`howto.exe` on Windows, `howto` on Linux/macOS). Download it from the link below, unpack and put somewhere in your `PATH` ([what's that?](https://gist.github.com/nex3/c395b2f8fd4b02068be37c961301caa7)), so you can run it from anyhwere on your computer.

[**Download**](https://github.com/nalgeon/howto/releases/latest)

**Note for macOS users**. macOS disables unsigned binaries and prevents the `howto` from running. To resolve this issue, remove the build from quarantine by running the following command in Terminal (replace `/path/to/folder` with an actual path to the folder containing the `howto` binary):

```text
xattr -d com.apple.quarantine /path/to/folder/howto
```

## Configuration

Howto is configured using environment variables. It can use cloud AIs or local Ollama models (coming soon).

Cloud AI providers charge for using their API, except for Gemini, which offers a free plan but may use your data in their products. Ollama is free without conditions but uses your machine's CPU or GPU resources.

Here's how to set up an AI provider:

### OpenAI

1. Get an API key from the [OpenAI Settings](https://platform.openai.com/account/api-keys).
2. Save the key to the `HOWTO_AI_TOKEN` environment variable.
3. Optionally set the `HOWTO_AI_MODEL` environment variable to the model name you want to use (default is `gpt-4o`).

### OpenAI-compatible provider

Anything like [OpenRouter](https://openrouter.ai/docs/), [Nebius](https://docs.nebius.com/studio/inference/api) or [Gemini](https://ai.google.dev/gemini-api/docs/openai):

1. Obtain an API endpoint from the documentation and save it to the `HOWTO_AI_URL` environment variable. Here are the endpoints for common providers:

-   OpenRouter: `https://openrouter.ai/api/v1/chat/completions`
-   Nebius: `https://api.studio.nebius.ai/v1/chat/completions`
-   Gemini: `https://generativelanguage.googleapis.com/v1beta/openai/chat/completions`

2. Get an API key from the provider and save it to the `HOWTO_AI_TOKEN` environment variable.
3. Set the `HOWTO_AI_MODEL` environment variable to the model name you want to use.

### Ollama (coming soon)

Ollama runs AI models locally on your machine. Here's how to set it up:

1. Download and install [Ollama](https://ollama.com/) for your operating system.
2. Set the [environment variables](https://github.com/ollama/ollama/blob/main/docs/faq.md#how-do-i-configure-ollama-server) to use less memory:

```text
OLLAMA_KEEP_ALIVE = 1h
OLLAMA_FLASH_ATTENTION = 1
```

3. Restart Ollama.
4. Download the AI model Gemma 2 (or another model of your choice):

```text
ollama pull gemma2:2b
```

5. Set the `HOWTO_AI_VENDOR` environment variable to `ollama`.
6. Set the `HOWTO_AI_MODEL` environment variable to `gemma2:2b` (or another model of your choice).

Gemma 2 is a lightweight model that uses about 1GB of memory and works quickly without a GPU.

### Other settings

-   `HOWTO_AI_TEMPERATURE`. Sampling temperature to use (between 0 and 2). Higher values make the output more random, while lower values make it more focused and predictable. Default: 0
-   `HOWTO_AI_TIMEOUT`. Timeout for AI API requests in seconds. Default: 30
-   `HOWTO_PROMPT`. The system prompt for the AI.

To see the system prompt and other settings, run `howto -v`.

## Usage

Describe your task to `howto`, and it will provide an answer:

```text
$ howto curl example.org but print only the headers
curl -I example.org

The `curl` command is used to transfer data from or to a server.
The `-I` option tells `curl` to fetch the HTTP headers only, without the body
content.
```

### Follow-ups

If you're not satisfied with an answer, refine it or ask a follow-up question by starting with `+`:

```text
$ howto a command that works kinda like diff but compares differently
comm file1 file2

The `comm` command compares two sorted files line by line and outputs three
columns: lines unique to the first file, lines unique to the second file, and
lines common to both files.

$ howto + yeah right i need only the intersection
comm -12 file1 file2

The `comm` command compares two sorted files line by line.
The `-12` option suppresses the first and second columns, showing only lines
common to both files (the intersection).
```

If you don't use `+`, howto will forget the previous conversation and treat your question as new.

### Run command

When satisfied with the suggested command, run `howto -run` to execute it without manually copying and pasting:

```text
$ howto curl example.org but print only the headers
curl -I example.org

The `curl` command is used to transfer data from or to a server.
The `-I` option tells `curl` to fetch the HTTP headers only, without the body
content.

$ howto -run
curl -I example.org

HTTP/1.1 200 OK
Content-Type: text/html
ETag: "84238dfc8092e5d9c0dac8ef93371a07:1736799080.121134"
Last-Modified: Mon, 13 Jan 2025 20:11:20 GMT
Cache-Control: max-age=2804
Date: Sun, 09 Feb 2025 12:54:51 GMT
Connection: keep-alive
```

That's it!

## License

Created by [Anton Zhiyanov](https://antonz.org/). Released under the MIT License.
