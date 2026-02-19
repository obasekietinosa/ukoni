package llm

import (
	"context"
	"ukoni/agent/internal/tools"
)

// Role constants
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Message represents a single message in the conversation.
type Message struct {
	Role       string
	Content    string
	ToolCallID string     // Required for RoleTool messages
	ToolCalls  []ToolCall // Populated for RoleAssistant messages that invoke tools
}

// Provider defines the interface for an LLM provider.
type Provider interface {
	// Generate generates a response from the LLM based on the conversation history and available tools.
	Generate(ctx context.Context, messages []Message, availableTools []tools.ToolDefinition) (Response, error)
	// Close cleans up resources used by the provider.
	Close() error
}

// Response represents the LLM's response.
type Response struct {
	Content   string
	ToolCalls []ToolCall
}

// ToolCall represents a request from the LLM to execute a tool.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string // JSON string of arguments
}
