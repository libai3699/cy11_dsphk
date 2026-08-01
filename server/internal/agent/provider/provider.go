package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model          string         `json:"model"`
	Messages       []Message      `json:"messages"`
	Temperature    float64        `json:"temperature,omitempty"`
	ResponseFormat map[string]any `json:"response_format,omitempty"`
}

type ChatResponse struct {
	Content string `json:"content"`
	Model   string `json:"model"`
}

type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewStepFunFromEnv() (*Client, bool) {
	apiKey := strings.TrimSpace(os.Getenv("STEP_API_KEY"))
	if apiKey == "" {
		return nil, false
	}

	baseURL := strings.TrimSpace(os.Getenv("STEP_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://api.stepfun.com/step_plan/v1"
	}

	model := strings.TrimSpace(os.Getenv("STEP_MODEL"))
	if model == "" {
		model = "step-3.7-flash"
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}, true
}

func (c *Client) Chat(ctx context.Context, messages []Message) (ChatResponse, error) {
	if c == nil || c.apiKey == "" {
		return ChatResponse{}, errors.New("stepfun api key is empty")
	}

	payload := ChatRequest{
		Model:          c.model,
		Messages:       messages,
		Temperature:    0.2,
		ResponseFormat: map[string]any{"type": "json_object"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return ChatResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ChatResponse{}, fmt.Errorf("stepfun status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return ChatResponse{}, err
	}
	if len(parsed.Choices) == 0 {
		return ChatResponse{}, errors.New("stepfun response has no choices")
	}

	return ChatResponse{
		Content: strings.TrimSpace(parsed.Choices[0].Message.Content),
		Model:   parsed.Model,
	}, nil
}

func ExtractJSONObject(raw string) string {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}
