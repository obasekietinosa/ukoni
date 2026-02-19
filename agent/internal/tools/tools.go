package tools

import (
	"context"
	"encoding/json"
)

// ToolDefinition represents the schema of a tool sent to the LLM.
type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"` // JSON Schema
}

// Tool represents an executable tool.
type Tool struct {
	Definition ToolDefinition
	Execute    func(ctx context.Context, arguments json.RawMessage) (string, error)
}
