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
	"strconv"
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
	agent   string
	apiKey  string
	baseURL string
	model   string
	timeout time.Duration
	client  *http.Client
}

type AgentConfig struct {
	Agent          string `json:"agent"`
	Provider       string `json:"provider"`
	BaseURL        string `json:"baseUrl"`
	Model          string `json:"model"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
	KeyConfigured  bool   `json:"keyConfigured"`
	Enabled        bool   `json:"enabled"`
}

func NewStepFunFromEnv() (*Client, bool) {
	return NewStepFunForAgent("")
}

func NewStepFunForAgent(agentName string) (*Client, bool) {
	apiKey := strings.TrimSpace(os.Getenv("STEP_API_KEY"))
	if apiKey == "" {
		return nil, false
	}

	config := StepFunConfigForAgent(agentName)
	if !config.Enabled {
		return nil, false
	}

	timeout := time.Duration(config.TimeoutSeconds) * time.Second

	return &Client{
		agent:   agentName,
		apiKey:  apiKey,
		baseURL: strings.TrimRight(config.BaseURL, "/"),
		model:   config.Model,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}, true
}

func StepFunConfigForAgent(agentName string) AgentConfig {
	baseURL := firstNonEmpty(agentEnv(agentName, "BASE_URL"), os.Getenv("STEP_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://api.stepfun.com/step_plan/v1"
	}

	model := firstNonEmpty(agentEnv(agentName, "MODEL"), os.Getenv("STEP_MODEL"))
	if model == "" {
		model = "step-3.7-flash"
	}

	timeoutSeconds := 120
	rawTimeout := firstNonEmpty(agentEnv(agentName, "TIMEOUT_SECONDS"), os.Getenv("STEP_TIMEOUT_SECONDS"))
	if rawTimeout != "" {
		if seconds, err := strconv.Atoi(rawTimeout); err == nil && seconds > 0 {
			timeoutSeconds = seconds
		}
	}

	enabled := true
	rawEnabled := strings.ToLower(firstNonEmpty(agentEnv(agentName, "ENABLED"), "true"))
	if rawEnabled == "0" || rawEnabled == "false" || rawEnabled == "off" || rawEnabled == "no" {
		enabled = false
	}

	return AgentConfig{
		Agent:          agentName,
		Provider:       "stepfun",
		BaseURL:        strings.TrimRight(baseURL, "/"),
		Model:          model,
		TimeoutSeconds: timeoutSeconds,
		KeyConfigured:  strings.TrimSpace(os.Getenv("STEP_API_KEY")) != "",
		Enabled:        enabled,
	}
}

func agentEnv(agentName string, suffix string) string {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(agentName), "-", "_"))
	if normalized == "" {
		return ""
	}
	return strings.TrimSpace(os.Getenv("AGENT_" + normalized + "_" + suffix))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
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
