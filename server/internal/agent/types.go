package agent

import "errors"

var ErrProviderNotConfigured = errors.New("agent provider is not configured")

type Result struct {
	Agent       string         `json:"agent"`
	Version     string         `json:"version"`
	Status      string         `json:"status"`
	Summary     string         `json:"summary"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Artifacts   map[string]any `json:"artifacts,omitempty"`
}

type RunOptions struct {
	TraceID    string `json:"traceId,omitempty"`
	OperatorID uint64 `json:"operatorId,omitempty"`
	DryRun     bool   `json:"dryRun,omitempty"`
}
