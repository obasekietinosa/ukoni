package openai

import (
	"context"
	"encoding/json"
	"ukoni/agent/internal/llm"
	"ukoni/agent/internal/tools"

	"github.com/sashabaranov/go-openai"
)

type Provider struct {
	client *openai.Client
	model  string
}

func New(apiKey string, model string) *Provider {
	if model == "" {
		model = openai.GPT4o
	}
	return &Provider{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (p *Provider) Generate(ctx context.Context, messages []llm.Message, availableTools []tools.ToolDefinition) (llm.Response, error) {
	// Convert messages
	chatMsgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msg := openai.ChatCompletionMessage{
			Content: m.Content,
		}
		switch m.Role {
		case llm.RoleSystem:
			msg.Role = openai.ChatMessageRoleSystem
		case llm.RoleUser:
			msg.Role = openai.ChatMessageRoleUser
		case llm.RoleAssistant:
			msg.Role = openai.ChatMessageRoleAssistant
			if len(m.ToolCalls) > 0 {
				// Convert tool calls
				tcs := make([]openai.ToolCall, len(m.ToolCalls))
				for j, tc := range m.ToolCalls {
					tcs[j] = openai.ToolCall{
						ID:   tc.ID,
						Type: openai.ToolTypeFunction,
						Function: openai.FunctionCall{
							Name:      tc.Name,
							Arguments: tc.Arguments,
						},
					}
				}
				msg.ToolCalls = tcs
			}
		case llm.RoleTool:
			msg.Role = openai.ChatMessageRoleTool
			msg.ToolCallID = m.ToolCallID
		}
		chatMsgs[i] = msg
	}

	// Convert tools
	var openAiTools []openai.Tool
	if len(availableTools) > 0 {
		openAiTools = make([]openai.Tool, len(availableTools))
		for i, t := range availableTools {
			// Parameters is json.RawMessage, but OpenAI expects a map/struct that marshals to JSON Schema.
			// We need to unmarshal the raw JSON into map[string]interface{}
			var params map[string]interface{}
			if len(t.Parameters) > 0 {
				if err := json.Unmarshal(t.Parameters, &params); err != nil {
					// If error unmarshalling, just pass nil or empty map
					params = make(map[string]interface{})
				}
			} else {
				params = make(map[string]interface{})
			}

			openAiTools[i] = openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  params,
				},
			}
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: chatMsgs,
		Tools:    openAiTools,
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return llm.Response{}, err
	}

	choice := resp.Choices[0]
	llmResp := llm.Response{
		Content: choice.Message.Content,
	}

	if len(choice.Message.ToolCalls) > 0 {
		tcs := make([]llm.ToolCall, len(choice.Message.ToolCalls))
		for i, tc := range choice.Message.ToolCalls {
			tcs[i] = llm.ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			}
		}
		llmResp.ToolCalls = tcs
	}

	return llmResp, nil
}

func (p *Provider) Close() error {
	return nil
}
