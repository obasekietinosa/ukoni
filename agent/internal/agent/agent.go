package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ukoni/agent/internal/llm"
	"ukoni/agent/internal/llm/gemini"
	"ukoni/agent/internal/llm/openai"
	"ukoni/agent/internal/tools"
	"ukoni/agent/pkg/client"
)

type Agent struct {
	APIBaseURL string
}

func New(apiBaseURL string) *Agent {
	return &Agent{APIBaseURL: apiBaseURL}
}

func (a *Agent) Run(ctx context.Context, prompt string, inventoryID string, token string) (string, error) {
	// Create client with auth
	c, err := client.NewClientWithResponses(a.APIBaseURL, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", token)
		return nil
	}))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	// Fetch Settings
	settingsResp, err := c.GetInventoriesIdSettingsWithResponse(ctx, inventoryID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch settings: %w", err)
	}
	if settingsResp.StatusCode() != 200 {
		return "", fmt.Errorf("inventory settings not found or accessible (status %d)", settingsResp.StatusCode())
	}

	settings := settingsResp.JSON200
	if settings == nil || settings.LlmApiKey == nil || *settings.LlmApiKey == "" {
		return "", fmt.Errorf("LLM API key not configured for this inventory. Please configure it in Inventory Settings.")
	}

	providerName := "openai"
	if settings.LlmProvider != nil {
		providerName = *settings.LlmProvider
	}

	var provider llm.Provider
	if providerName == "gemini" {
		p, err := gemini.New(ctx, *settings.LlmApiKey, "")
		if err != nil {
			return "", fmt.Errorf("failed to initialize gemini: %w", err)
		}
		provider = p
	} else {
		provider = openai.New(*settings.LlmApiKey, "")
	}
	defer provider.Close()

	// Initialize ToolSet
	toolSet := tools.NewToolSet(c)
	availableTools := toolSet.GetTools()
	var toolDefs []tools.ToolDefinition
	for _, t := range availableTools {
		toolDefs = append(toolDefs, t.Definition)
	}

	// Messages
	messages := []llm.Message{
		{Role: llm.RoleSystem, Content: "You are a helpful household assistant for Ukoni. You can manage inventory and shopping lists. Current Inventory ID: " + inventoryID},
		{Role: llm.RoleUser, Content: prompt},
	}

	// Loop
	maxTurns := 5
	for i := 0; i < maxTurns; i++ {
		resp, err := provider.Generate(ctx, messages, toolDefs)
		if err != nil {
			return "", fmt.Errorf("LLM generation error: %w", err)
		}

		// Add assistant message
		msg := llm.Message{
			Role:      llm.RoleAssistant,
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		}
		messages = append(messages, msg)

		if len(resp.ToolCalls) == 0 {
			return resp.Content, nil
		}

		// Execute tools
		for _, tc := range resp.ToolCalls {
			// Find tool
			var tool *tools.Tool
			// We need to iterate by value if we want to take address, but here we just copy
			for _, t := range availableTools {
				if t.Definition.Name == tc.Name {
					tool = &t // Be careful with loop variable address in older Go, but safe in 1.22+
					break
				}
			}

			resultContent := ""
			if tool != nil {
				// We need to use 'tool' variable which is *Tool
				// Wait, in range loop 't' is a copy. So taking address of t is taking address of copy.
				// But we just need to call Execute on it.
				// Also passing 't' (by value) is fine if Execute is safe.
				// Tool struct has 'Execute' field which is a function.
				// So:
				res, err := tool.Execute(ctx, json.RawMessage(tc.Arguments))
				if err != nil {
					resultContent = fmt.Sprintf("Error executing tool %s: %v", tc.Name, err)
				} else {
					resultContent = res
				}
			} else {
				resultContent = fmt.Sprintf("Error: Tool %s not found", tc.Name)
			}

			// Add tool response
			messages = append(messages, llm.Message{
				Role:       llm.RoleTool,
				Content:    resultContent,
				ToolCallID: tc.ID,
			})
		}
	}

	return messages[len(messages)-1].Content, nil
}
