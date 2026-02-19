package gemini

import (
	"encoding/json"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

func TestJsonSchemaToSchema(t *testing.T) {
	rawSchema := `{
		"type": "object",
		"properties": {
			"inventory_id": {
				"type": "string",
				"description": "The ID of the inventory."
			},
			"products": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"name": {
							"type": "string"
						}
					},
					"required": ["name"]
				}
			}
		},
		"required": ["inventory_id", "products"]
	}`

	schema, err := parseJSONSchema(json.RawMessage(rawSchema))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	if schema.Type != genai.TypeObject {
		t.Errorf("Expected TypeObject, got %v", schema.Type)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}

	invID, ok := schema.Properties["inventory_id"]
	if !ok {
		t.Fatal("inventory_id property missing")
	}
	if invID.Type != genai.TypeString {
		t.Errorf("Expected inventory_id to be String, got %v", invID.Type)
	}
	if invID.Description != "The ID of the inventory." {
		t.Errorf("Expected description match, got %s", invID.Description)
	}

	products, ok := schema.Properties["products"]
	if !ok {
		t.Fatal("products property missing")
	}
	if products.Type != genai.TypeArray {
		t.Errorf("Expected products to be Array, got %v", products.Type)
	}
	if products.Items == nil {
		t.Fatal("products items missing")
	}
	if products.Items.Type != genai.TypeObject {
		t.Errorf("Expected products items to be Object, got %v", products.Items.Type)
	}
}
