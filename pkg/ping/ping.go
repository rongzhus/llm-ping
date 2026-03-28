package ping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultMaxRetries    = 4
	DefaultRetryInterval = 1 * time.Second
	DefaultTimeout       = 30 * time.Second
)

type PingResult struct {
	Success     bool
	Message     string
	Response    string
	Latency     time.Duration
	RetryCount  int
	Error       error
}

type ModelRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ModelResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func PingModel(host, port, endpoint, apiKey, modelName string) PingResult {
	return PingModelWithRetry(host, port, endpoint, apiKey, modelName, DefaultMaxRetries, DefaultRetryInterval)
}

func PingModelWithRetry(host, port, endpoint, apiKey, modelName string, maxRetries int, retryInterval time.Duration) PingResult {
	url := fmt.Sprintf("%s:%s%s", host, port, endpoint)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result := pingOnce(url, apiKey, modelName)
		
		if result.Success {
			return result
		}

		if attempt < maxRetries {
			time.Sleep(retryInterval)
		}

		result.RetryCount = attempt + 1
		if attempt == maxRetries {
			return result
		}
	}

	return PingResult{
		Success:    false,
		Message:    fmt.Sprintf("Failed after %d attempts", maxRetries+1),
		RetryCount: maxRetries,
		Error:      fmt.Errorf("max retries exceeded"),
	}
}

func pingOnce(url, apiKey, modelName string) PingResult {
	startTime := time.Now()

	reqBody := ModelRequest{
		Model: modelName,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: "Respond with just 'OK' and nothing else",
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return PingResult{
			Success: false,
			Message: "Failed to marshal request",
			Error:   err,
		}
	}

	client := &http.Client{
		Timeout: DefaultTimeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return PingResult{
			Success: false,
			Message: "Failed to create request",
			Error:   err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return PingResult{
			Success: false,
			Message: fmt.Sprintf("Request failed: %v", err),
			Error:   err,
		}
	}
	defer resp.Body.Close()

	latency := time.Since(startTime)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PingResult{
			Success: false,
			Message: "Failed to read response",
			Latency: latency,
			Error:   err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return PingResult{
			Success: false,
			Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
			Latency: latency,
			Error:   fmt.Errorf("HTTP %d", resp.StatusCode),
		}
	}

	var modelResp ModelResponse
	if err := json.Unmarshal(body, &modelResp); err != nil {
		return PingResult{
			Success: false,
			Message: "Failed to parse response",
			Latency: latency,
			Error:   err,
		}
	}

	if modelResp.Error != nil {
		return PingResult{
			Success: false,
			Message: fmt.Sprintf("API Error: %s", modelResp.Error.Message),
			Latency: latency,
			Error:   fmt.Errorf("API error: %s", modelResp.Error.Message),
		}
	}

	if len(modelResp.Choices) == 0 {
		return PingResult{
			Success: false,
			Message: "No choices in response",
			Latency: latency,
			Error:   fmt.Errorf("no choices in response"),
		}
	}

	response := modelResp.Choices[0].Message.Content
	if response == "OK" {
		return PingResult{
			Success:  true,
			Message:  "OK",
			Response: response,
			Latency:  latency,
		}
	}

	return PingResult{
		Success:  false,
		Message:  fmt.Sprintf("Unexpected response: %s", response),
		Response: response,
		Latency:  latency,
		Error:    fmt.Errorf("unexpected response"),
	}
}
