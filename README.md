# llm-ping

A simple command-line tool to ping and test LLM (Large Language Model) services via REST API. Similar to network ping, it allows you to check if your LLM service is responding correctly.

## Features

- **Ping LLM Services**: Test if your LLM API is responding correctly
- **Auto Retry**: Automatic retry mechanism with configurable attempts and intervals
- **Configuration Management**: Save and manage multiple model configurations
- **Simple Output**: Returns "OK" on success, detailed error messages on failure
- **Cross-platform**: Supports Windows, macOS, and Linux

## Installation

### From Source

```bash
git clone https://github.com/your-username/llm-ping.git
cd llm-ping
go build -o llm-ping main.go
```

### Using Make (Recommended)

```bash
make build
```

This will build the binary for your current platform.

## Usage

### Configure a Model

Save a model configuration for easy access:

```bash
llm-ping set mymodel --host http://localhost --port 8080 --api-key sk-xxx
```

Options:
- `--host`: Model host URL (required)
- `--port`: Model port (required)
- `--endpoint`: API endpoint (default: `/v1/chat/completions`)
- `--api-key`: API key for authentication (optional)

### Ping a Model

#### Using saved configuration:

```bash
llm-ping ping mymodel
```

#### Using direct parameters:

```bash
llm-ping ping --host http://localhost --port 8080 --api-key sk-xxx
```

#### With custom retry settings:

```bash
llm-ping ping mymodel --retries 5 --interval 2
```

Options:
- `--host`: Model host URL (required if not using saved config)
- `--port`: Model port (required if not using saved config)
- `--endpoint`: API endpoint (default: `/v1/chat/completions`)
- `--api-key`: API key for authentication
- `--model`: Model name (default: `gpt-3.5-turbo`)
- `--retries`: Number of retry attempts (default: 4)
- `--interval`: Retry interval in seconds (default: 1)

### List Configured Models

```bash
llm-ping list
```

This will display all saved model configurations.

### Remove a Model Configuration

```bash
llm-ping remove mymodel
```

## Output Examples

### Successful Ping

```bash
$ llm-ping ping mymodel
OK (245ms)
```

### Failed Ping

```bash
$ llm-ping ping mymodel
Error: HTTP 503: Service Unavailable (retried 4 times)
```

### Connection Error

```bash
$ llm-ping ping --host http://localhost --port 9999
Error: Request failed: dial tcp: connect: connection refused (retried 4 times)
```

## Configuration File

Configuration files are stored in `~/.llm-ping/llm-ping.json` (Unix) or `%USERPROFILE%\.llm-ping\llm-ping.json` (Windows).

Example configuration:

```json
{
  "models": [
    {
      "name": "openai",
      "host": "https://api.openai.com",
      "port": "443",
      "endpoint": "/v1/chat/completions",
      "apiKey": "sk-xxxxxxxxxxxxxxxx"
    },
    {
      "name": "local-model",
      "host": "http://localhost",
      "port": "8080",
      "endpoint": "/v1/chat/completions",
      "apiKey": ""
    }
  ]
}
```

## Building for Different Platforms

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platform
make build-linux
make build-windows
make build-macos
```

### Manual Build

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o llm-ping-linux main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o llm-ping.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o llm-ping-mac main.go
```

## API Compatibility

This tool is compatible with OpenAI-compatible APIs, including:
- OpenAI API
- Azure OpenAI
- Local LLM servers (e.g., vLLM, Ollama with OpenAI-compatible endpoints)
- Other OpenAI-compatible services

## Requirements

- Go 1.21 or higher
- No external dependencies beyond Go standard library

## Project Structure

```
llm-ping/
├── main.go              # Main entry point and CLI handling
├── go.mod              # Go module definition
├── pkg/
│   ├── config/
│   │   └── config.go   # Configuration management
│   └── ping/
│       └── ping.go     # Core ping functionality
├── Makefile            # Build scripts
└── README.md           # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
