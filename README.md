# Chat CLI

A modern, interactive command-line interface for chatting with Large Language Models (LLMs). Built with Go and featuring both basic and Terminal User Interface (TUI) modes for seamless AI conversations.

## 🚀 Features

- **Multiple LLM Providers**: Currently supports Claude/Anthropic with plans for OpenAI GPT and local models
- **Dual Interface Modes**: 
  - **TUI Mode** (Default): Rich terminal interface with markdown rendering and streaming responses
  - **Basic Mode**: Simple readline-based interface for lightweight usage
- **Streaming Support**: Real-time response streaming for natural conversation flow
- **Session Management**: Maintains conversation context throughout your chat session
- **Flexible Configuration**: Environment variables, command-line flags, and config file support
- **Cross-Platform**: Works on Linux, macOS, and Windows

## 📦 Installation

### Prerequisites
- Go 1.21 or later
- An Anthropic API key

### From Source
```bash
git clone https://github.com/donaldknoller/chat-cli.git
cd chat-cli
go build -o chat-cli ./cmd/main.go
```

### Using Go Install
```bash
go install github.com/donaldknoller/chat-cli/cmd@latest
```

## ⚙️ Configuration

### API Key Setup
You can provide your Anthropic API key in several ways:

1. **Environment Variable** (Recommended):
   ```bash
   export CHAT_CLI_LLM_API_KEY="your-anthropic-api-key"
   ```

2. **Command Line Flag**:
   ```bash
   ./chat-cli --key "your-anthropic-api-key"
   ```

3. **Key File** (Default location: `~/.config/.chat_cli/key`):
   ```bash
   mkdir -p ~/.config/.chat_cli
   echo "your-anthropic-api-key" > ~/.config/.chat_cli/key
   ```

### Configuration Options

| Environment Variable | CLI Flag | Default Value | Description |
|---------------------|----------|---------------|-------------|
| `CHAT_CLI_LLM_API_KEY` | `--key` | - | Anthropic API key |
| `CHAT_CLI_LLM_API_HOST` | `--host` | `https://api.anthropic.com/v1/messages` | API endpoint |
| `CHAT_CLI_LLM_API_MODEL` | - | `claude-3-opus-20240229` | Claude model to use |
| `CHAT_CLI_SINGLE_RUN_MODE` | - | `false` | Enable single-run mode |
| `CHAT_CLI_DEBUG` | - | `false` | Enable debug logging |

## 🎯 Usage

### TUI Mode (Default)
The rich terminal interface provides the best chat experience:

```bash
./chat-cli
```

**TUI Controls:**
- **Enter**: Add new lines to your message
- **Ctrl+S**: Send message to the AI
- **Ctrl+C**: Exit the application
- **Arrow Keys**: Navigate through chat history

### Basic Mode
For simple, lightweight interactions:

```bash
./chat-cli --single-run-mode
```

**Basic Mode Features:**
- Simple `>` prompt for input
- Direct message sending (no multiline support)
- Suitable for scripting and automation

### Example Usage

```bash
# Start interactive chat
./chat-cli

# Use specific API key
./chat-cli --key "sk-ant-api03-..."

# Use different model via environment
export CHAT_CLI_LLM_API_MODEL="claude-3-sonnet-20240229"
./chat-cli
```

## 🏗️ Development

### Project Structure
```
chat-cli/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── anthropic/           # Anthropic API client
│   │   ├── client.go        # HTTP client
│   │   ├── client_stream.go # Streaming client
│   │   ├── request.go       # Request structures
│   │   ├── response.go      # Response handling
│   │   └── session.go       # Session management
│   ├── common_llm/          # Common LLM interfaces
│   │   ├── models.go        # Shared models
│   │   └── session.go       # Session interface
│   ├── config/              # Configuration management
│   │   └── config.go        # Config loading and defaults
│   └── interactive/         # TUI implementation
│       ├── entrypoint.go    # TUI entry point
│       └── tea.go           # Bubble Tea UI components
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── LICENSE                  # Apache 2.0 License
└── README.md               # This file
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/donaldknoller/chat-cli.git
cd chat-cli

# Download dependencies
go mod download

# Build the application
go build -o chat-cli ./cmd/main.go

# Run tests
go test ./...
```

### Dependencies
- **[Cobra](https://github.com/spf13/cobra)**: CLI framework
- **[Viper](https://github.com/spf13/viper)**: Configuration management
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: TUI framework
- **[Glamour](https://github.com/charmbracelet/glamour)**: Markdown rendering
- **[Readline](https://github.com/chzyer/readline)**: Line editing for basic mode

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🗺️ Roadmap

- [ ] OpenAI GPT support
- [ ] Local model support (Ollama, etc.)
- [ ] Conversation history persistence
- [ ] Custom prompt templates
- [ ] Plugin system
- [ ] Web interface
- [ ] Docker containerization

## 💡 Tips

- Use the TUI mode for the best experience with syntax highlighting and markdown rendering
- Set your API key as an environment variable for convenience
- The application maintains conversation context within a session
- Use Ctrl+C to exit cleanly from either mode