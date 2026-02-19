package gemini

import (
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// jsonSchemaToSchema converts a JSON schema (map[string]interface{}) to a *genai.Schema.
func jsonSchemaToSchema(schema map[string]interface{}) *genai.Schema {
	if schema == nil {
		return nil
	}

	t, _ := schema["type"].(string)

	// Map type string to genai.Type
	var genaiType genai.Type
	switch t {
	case "object":
		genaiType = genai.TypeObject
	case "array":
		genaiType = genai.TypeArray
	case "string":
		genaiType = genai.TypeString
	case "number":
		genaiType = genai.TypeNumber
	case "integer":
		genaiType = genai.TypeInteger
	case "boolean":
		genaiType = genai.TypeBoolean
	default:
		// Default or fallback
		genaiType = genai.TypeString
	}

	res := &genai.Schema{
		Type: genaiType,
	}

	if desc, ok := schema["description"].(string); ok {
		res.Description = desc
	}

	// Handle properties for objects
	if props, ok := schema["properties"].(map[string]interface{}); ok {
		res.Properties = make(map[string]*genai.Schema)
		for k, v := range props {
			if vMap, ok := v.(map[string]interface{}); ok {
				res.Properties[k] = jsonSchemaToSchema(vMap)
			}
		}
	}

	// Handle required fields
	if req, ok := schema["required"].([]interface{}); ok {
		for _, r := range req {
			if rStr, ok := r.(string); ok {
				res.Required = append(res.Required, rStr)
			}
		}
	}

	// Handle items for arrays
	if items, ok := schema["items"].(map[string]interface{}); ok {
		res.Items = jsonSchemaToSchema(items)
	}

	// Handle enum
	if enum, ok := schema["enum"].([]interface{}); ok {
		for _, e := range enum {
			if eStr, ok := e.(string); ok {
				res.Enum = append(res.Enum, eStr)
			}
		}
	}

	return res
}

func parseJSONSchema(raw json.RawMessage) (*genai.Schema, error) {
	var schemaMap map[string]interface{}
	if err := json.Unmarshal(raw, &schemaMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	return jsonSchemaToSchema(schemaMap), nil
}
