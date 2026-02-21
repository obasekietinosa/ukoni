package tools

import (
	"encoding/json"
	"testing"
	"ukoni/agent/pkg/client"
)

func TestPlansTools_Definition(t *testing.T) {
	ts := NewToolSet(&client.ClientWithResponses{})

	tests := []struct {
		name     string
		toolFunc func() Tool
		toolName string
	}{
		{"CreatePlanTool", ts.CreatePlanTool, "create_plan"},
		{"AddPlanItemTool", ts.AddPlanItemTool, "add_plan_item"},
		{"CreatePlanGroupTool", ts.CreatePlanGroupTool, "create_plan_group"},
		{"AddPlanToGroupTool", ts.AddPlanToGroupTool, "add_plan_to_group"},
		{"CreateShoppingListFromPlanTool", ts.CreateShoppingListFromPlanTool, "create_shopping_list_from_plan"},
		{"CreateShoppingListFromPlanGroupTool", ts.CreateShoppingListFromPlanGroupTool, "create_shopping_list_from_plan_group"},
		{"ListPlansTool", ts.ListPlansTool, "list_plans"},
		{"ListPlanGroupsTool", ts.ListPlanGroupsTool, "list_plan_groups"},
		{"GetPlanTool", ts.GetPlanTool, "get_plan"},
		{"GetPlanGroupTool", ts.GetPlanGroupTool, "get_plan_group"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := tt.toolFunc()
			if tool.Definition.Name != tt.toolName {
				t.Errorf("expected tool name '%s', got '%s'", tt.toolName, tool.Definition.Name)
			}

			// Verify parameters schema is valid JSON
			var schema map[string]interface{}
			if err := json.Unmarshal(tool.Definition.Parameters, &schema); err != nil {
				t.Errorf("invalid parameters schema for tool '%s': %v", tt.toolName, err)
			}
		})
	}
}
