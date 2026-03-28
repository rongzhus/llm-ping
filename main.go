package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/llm-ping/pkg/config"
	"github.com/llm-ping/pkg/ping"
)

var (
	host     string
	port     string
	endpoint string
	apiKey   string
	model    string
	retries  int
	interval int
)

func init() {
	flag.StringVar(&host, "host", "", "Model host URL (e.g., http://localhost)")
	flag.StringVar(&port, "port", "", "Model port (e.g., 8080)")
	flag.StringVar(&endpoint, "endpoint", "/v1/chat/completions", "API endpoint")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication")
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "Model name")
	flag.IntVar(&retries, "retries", 4, "Number of retry attempts")
	flag.IntVar(&interval, "interval", 1, "Retry interval in seconds")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "llm-ping - A simple tool to ping LLM models via REST API\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  llm-ping <command> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  ping      Ping a configured model or specified host\n")
		fmt.Fprintf(os.Stderr, "  set       Configure a new model\n")
		fmt.Fprintf(os.Stderr, "  config    Configure a new model (alias for set)\n")
		fmt.Fprintf(os.Stderr, "  list      List all configured models\n")
		fmt.Fprintf(os.Stderr, "  remove    Remove a configured model\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  llm-ping set mymodel --host http://localhost --port 8080 --api-key sk-xxx\n")
		fmt.Fprintf(os.Stderr, "  llm-ping ping mymodel\n")
		fmt.Fprintf(os.Stderr, "  llm-ping ping --host http://localhost --port 8080\n")
		fmt.Fprintf(os.Stderr, "  llm-ping list\n")
		fmt.Fprintf(os.Stderr, "  llm-ping remove mymodel\n")
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "set" || command == "config" {
		if len(os.Args) < 3 {
			fmt.Println("Error: Model name is required")
			os.Exit(1)
		}
		modelName := os.Args[2]
		args := os.Args[3:]
		flag.CommandLine.Parse(args)
		handleSetWithModelName(modelName)
		return
	}

	flag.CommandLine.Parse(os.Args[2:])

	switch command {
	case "ping":
		handlePing()
	case "set", "config":
		handleSet()
	case "list":
		handleList()
	case "remove":
		handleRemove()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

func handlePing() {
	var pingHost, pingPort, pingEndpoint, pingAPIKey, pingModel string

	if flag.NArg() > 0 {
		modelName := flag.Arg(0)
		cfg, err := config.GetModel(modelName)
		if err != nil {
			fmt.Printf("Error: Failed to load config: %v\n", err)
			os.Exit(1)
		}
		if cfg == nil {
			fmt.Printf("Error: Model '%s' not found in config\n", modelName)
			os.Exit(1)
		}
		pingHost = cfg.Host
		pingPort = cfg.Port
		pingEndpoint = cfg.Endpoint
		pingAPIKey = cfg.APIKey
		if cfg.ModelName != "" {
			pingModel = cfg.ModelName
		} else {
			pingModel = cfg.Name
		}
	} else {
		if host == "" {
			fmt.Println("Error: --host is required when not using a configured model")
			os.Exit(1)
		}
		pingHost = host
		pingPort = port
		pingEndpoint = endpoint
		pingAPIKey = apiKey
		pingModel = model
	}

	if pingPort == "" {
		pingPort = "443"
	}

	result := ping.PingModelWithRetry(
		pingHost,
		pingPort,
		pingEndpoint,
		pingAPIKey,
		pingModel,
		retries,
		1000*time.Duration(interval),
	)

	if result.Success {
		fmt.Printf("OK (%dms)\n", result.Latency.Milliseconds())
	} else {
		fmt.Printf("Error: %s", result.Message)
		if result.Error != nil {
			fmt.Printf(" (%v)", result.Error)
		}
		if result.RetryCount > 0 {
			fmt.Printf(" (retried %d times)", result.RetryCount)
		}
		fmt.Println()
		os.Exit(1)
	}
}

func handleSet() {
	if flag.NArg() == 0 {
		fmt.Println("Error: Model name is required")
		os.Exit(1)
	}

	modelName := flag.Arg(0)
	handleSetWithModelName(modelName)
}

func handleSetWithModelName(modelName string) {
	if host == "" {
		fmt.Println("Error: --host is required")
		os.Exit(1)
	}

	if port == "" {
		fmt.Println("Error: --port is required")
		os.Exit(1)
	}

	cfg := config.ModelConfig{
		Name:      modelName,
		Host:      host,
		Port:      port,
		Endpoint:  endpoint,
		APIKey:    apiKey,
		ModelName: model,
	}

	if err := config.AddModel(cfg); err != nil {
		fmt.Printf("Error: Failed to save config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model '%s' configured successfully\n", modelName)
}

func handleList() {
	models, err := config.ListModels()
	if err != nil {
		fmt.Printf("Error: Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if len(models) == 0 {
		fmt.Println("No models configured")
		return
	}

	fmt.Println("Configured models:")
	fmt.Println(strings.Repeat("-", 80))
	for _, m := range models {
		fmt.Printf("Name:     %s\n", m.Name)
		if m.ModelName != "" {
			fmt.Printf("  Model:  %s\n", m.ModelName)
		}
		fmt.Printf("  Host:   %s\n", m.Host)
		fmt.Printf("  Port:   %s\n", m.Port)
		fmt.Printf("  Endpoint: %s\n", m.Endpoint)
		if m.APIKey != "" {
			fmt.Printf("  API Key:  %s...\n", m.APIKey[:min(10, len(m.APIKey))])
		}
		fmt.Println()
	}
}

func handleRemove() {
	if flag.NArg() == 0 {
		fmt.Println("Error: Model name is required")
		os.Exit(1)
	}

	modelName := flag.Arg(0)

	if err := config.RemoveModel(modelName); err != nil {
		fmt.Printf("Error: Failed to remove model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model '%s' removed successfully\n", modelName)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
