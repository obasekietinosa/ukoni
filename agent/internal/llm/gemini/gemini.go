package gemini

import (
	"context"
	"fmt"
	"strings"
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
	// Tool support not implemented for Gemini in MVP

	// Construct history as a single prompt
	var promptBuilder strings.Builder
	for _, m := range messages {
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", m.Role, m.Content))
	}

	// Send request
	resp, err := model.GenerateContent(ctx, genai.Text(promptBuilder.String()))
	if err != nil {
		return llm.Response{}, err
	}

	// Parse response
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return llm.Response{Content: ""}, nil
	}

	part := resp.Candidates[0].Content.Parts[0]
	if txt, ok := part.(genai.Text); ok {
		return llm.Response{Content: string(txt)}, nil
	}

	return llm.Response{Content: "Received non-text response from Gemini"}, nil
}

func (p *Provider) Close() error {
	return p.client.Close()
}
