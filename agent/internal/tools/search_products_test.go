package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"ukoni/agent/pkg/client"
)

func TestSearchCanonicalProductsTool_Definition(t *testing.T) {
	ts := NewToolSet(&client.ClientWithResponses{})
	tool := ts.SearchCanonicalProductsTool()

	if tool.Definition.Name != "search_canonical_products" {
		t.Errorf("expected tool name 'search_canonical_products', got '%s'", tool.Definition.Name)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(tool.Definition.Parameters, &schema); err != nil {
		t.Errorf("invalid parameters schema: %v", err)
	}
}

func TestSearchCanonicalProductsTool_InputParsing(t *testing.T) {
	ts := NewToolSet(&client.ClientWithResponses{})
	tool := ts.SearchCanonicalProductsTool()

	// Provide invalid arguments to ensure unmarshal catches it
	invalidArgs := json.RawMessage(`{"inventory_id": 123}`) // integer instead of string
	_, err := tool.Execute(context.Background(), invalidArgs)
	if err == nil {
		t.Errorf("expected error with invalid arguments, got nil")
	} else if !strings.Contains(err.Error(), "invalid arguments") {
		t.Errorf("expected 'invalid arguments' error, got: %v", err)
	}
}
