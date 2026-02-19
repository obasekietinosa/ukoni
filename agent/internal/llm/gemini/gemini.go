package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"ukoni/agent/internal/llm"
	"ukoni/agent/internal/tools"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Provider struct {
	client *genai.Client
	model  string
}

func New(ctx context.Context, apiKey string, model string) (*Provider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Provider{
		client: client,
		model:  model,
	}, nil
}

func (p *Provider) Generate(ctx context.Context, messages []llm.Message, availableTools []tools.ToolDefinition) (llm.Response, error) {
	model := p.client.GenerativeModel(p.model)

	// Configure tools
	if len(availableTools) > 0 {
		var genaiTools []*genai.Tool
		for _, t := range availableTools {
			schema, err := parseJSONSchema(t.Parameters)
			if err != nil {
				fmt.Printf("Error parsing schema for tool %s: %v\n", t.Name, err)
				continue
			}

			genaiTools = append(genaiTools, &genai.Tool{
				FunctionDeclarations: []*genai.FunctionDeclaration{{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  schema,
				}},
			})
		}
		model.Tools = genaiTools
	}

	// Convert history
	var history []*genai.Content
	var systemInstruction *genai.Content

	for i, m := range messages {
		switch m.Role {
		case llm.RoleSystem:
			systemInstruction = &genai.Content{
				Parts: []genai.Part{genai.Text(m.Content)},
			}
		case llm.RoleUser:
			history = append(history, &genai.Content{
				Role: "user",
				Parts: []genai.Part{genai.Text(m.Content)},
			})
		case llm.RoleAssistant:
			parts := []genai.Part{}
			if m.Content != "" {
				parts = append(parts, genai.Text(m.Content))
			}
			for _, tc := range m.ToolCalls {
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
					args = make(map[string]interface{})
				}
				parts = append(parts, genai.FunctionCall{
					Name: tc.Name,
					Args: args,
				})
			}
			history = append(history, &genai.Content{
				Role: "model",
				Parts: parts,
			})
		case llm.RoleTool:
			funcName := "unknown"
			found := false
			// Search backwards for the tool call ID
			for j := i - 1; j >= 0; j-- {
				if messages[j].Role == llm.RoleAssistant {
					for _, tc := range messages[j].ToolCalls {
						if tc.ID == m.ToolCallID {
							funcName = tc.Name
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}

			var responseData map[string]interface{}
			if err := json.Unmarshal([]byte(m.Content), &responseData); err != nil {
				responseData = map[string]interface{}{
					"result": m.Content,
				}
			}

			history = append(history, &genai.Content{
				Role: "function",
				Parts: []genai.Part{genai.FunctionResponse{
					Name:     funcName,
					Response: responseData,
				}},
			})
		}
	}

	if systemInstruction != nil {
		model.SystemInstruction = systemInstruction
	}

	// We'll create a session but we need to manage history carefully.
	// We want to send the *last* message in the history list as the new message.
	if len(history) == 0 {
		return llm.Response{}, fmt.Errorf("no messages to send")
	}

	lastMsg := history[len(history)-1]

	// Start Chat
	session := model.StartChat()

	// Set history to everything *except* the last message
	if len(history) > 1 {
		session.History = history[:len(history)-1]
	}

	// Send the last message
	resp, err := session.SendMessage(ctx, lastMsg.Parts...)
	if err != nil {
		return llm.Response{}, err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return llm.Response{Content: ""}, nil
	}

	content := ""
	var toolCalls []llm.ToolCall
	callCount := 0

	for _, part := range resp.Candidates[0].Content.Parts {
		switch p := part.(type) {
		case genai.Text:
			content += string(p)
		case genai.FunctionCall:
			callCount++
			argsBytes, _ := json.Marshal(p.Args)
			toolCalls = append(toolCalls, llm.ToolCall{
				ID:        fmt.Sprintf("call_%s_%d", p.Name, callCount),
				Name:      p.Name,
				Arguments: string(argsBytes),
			})
		}
	}

	return llm.Response{
		Content:   content,
		ToolCalls: toolCalls,
	}, nil
}

func (p *Provider) Close() error {
	return p.client.Close()
}
