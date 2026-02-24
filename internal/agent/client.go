package agent

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

// Client wraps the Anthropic API client
type Client struct {
    apiKey      string
    model       string
    temperature float64
    httpClient  *http.Client
}

// NewClient creates a new agent client
func NewClient(model string, temperature float64) (*Client, error) {
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
    }

    return &Client{
        apiKey:      apiKey,
        model:       model,
        temperature: temperature,
        httpClient:  &http.Client{},
    }, nil
}

type messageRequest struct {
    Model       string    `json:"model"`
    MaxTokens   int       `json:"max_tokens"`
    Temperature float64   `json:"temperature"`
    System      string    `json:"system"`
    Messages    []message `json:"messages"`
}

type message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type messageResponse struct {
    ID      string         `json:"id"`
    Type    string         `json:"type"`
    Role    string         `json:"role"`
    Content []contentBlock `json:"content"`
    Model   string         `json:"model"`
    Usage   usage          `json:"usage"`
}

type contentBlock struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

type usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}

type errorResponse struct {
    Type  string      `json:"type"`
    Error errorDetail `json:"error"`
}

type errorDetail struct {
    Type    string `json:"type"`
    Message string `json:"message"`
}

// GenerateCompletion sends a message to Claude and returns the response
func (c *Client) GenerateCompletion(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
    reqBody := messageRequest{
        Model:       c.model,
        MaxTokens:   4096,
        Temperature: c.temperature,
        System:      systemPrompt,
        Messages: []message{
            {
                Role:    "user",
                Content: userPrompt,
            },
        },
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", c.apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        var errResp errorResponse
        if err := json.Unmarshal(body, &errResp); err != nil {
            return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
        }
        return "", fmt.Errorf("API error: %s - %s", errResp.Error.Type, errResp.Error.Message)
    }

    var msgResp messageResponse
    if err := json.Unmarshal(body, &msgResp); err != nil {
        return "", fmt.Errorf("failed to unmarshal response: %w", err)
    }

    if len(msgResp.Content) == 0 {
        return "", fmt.Errorf("empty response from Claude")
    }

    if msgResp.Content[0].Type != "text" {
        return "", fmt.Errorf("unexpected content type: %s", msgResp.Content[0].Type)
    }

    return msgResp.Content[0].Text, nil
}
