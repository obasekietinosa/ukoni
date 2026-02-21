package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"ukoni/agent/pkg/client"
)

func (ts *ToolSet) CreatePlanTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "create_plan",
			Description: "Create a new plan.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					},
					"title": {
						"type": "string",
						"description": "The title of the plan."
					},
					"description": {
						"type": "string",
						"description": "The description of the plan."
					}
				},
				"required": ["inventory_id", "title"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersCreatePlanRequest{
				Title:       &input.Title,
				Description: &input.Description,
			}
			resp, err := ts.Client.PostInventoriesIdPlansWithResponse(ctx, input.InventoryID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON201)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) AddPlanItemTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "add_plan_item",
			Description: "Add an item to a plan.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"plan_id": {
						"type": "string",
						"description": "The ID of the plan."
					},
					"target_type": {
						"type": "string",
						"enum": ["canonical_product", "product_variant", "product"],
						"description": "Type of item to add."
					},
					"target_id": {
						"type": "string",
						"description": "ID of the product/variant."
					},
					"quantity": {
						"type": "number",
						"description": "Quantity."
					},
					"unit": {
						"type": "string",
						"description": "Unit (e.g., kg, pcs)."
					},
					"note": {
						"type": "string",
						"description": "Note."
					}
				},
				"required": ["plan_id", "target_type", "target_id", "quantity"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				PlanID     string  `json:"plan_id"`
				TargetType string  `json:"target_type"`
				TargetID   string  `json:"target_id"`
				Quantity   float32 `json:"quantity"`
				Unit       string  `json:"unit"`
				Note       string  `json:"note"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersPlanItemRequest{
				TargetType: &input.TargetType,
				TargetId:   &input.TargetID,
				Quantity:   &input.Quantity,
				Unit:       &input.Unit,
				Note:       &input.Note,
			}
			resp, err := ts.Client.PostPlansIdItemsWithResponse(ctx, input.PlanID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON201)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) CreatePlanGroupTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "create_plan_group",
			Description: "Create a new plan group.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					},
					"title": {
						"type": "string",
						"description": "The title of the plan group."
					},
					"description": {
						"type": "string",
						"description": "The description of the plan group."
					}
				},
				"required": ["inventory_id", "title"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersCreatePlanGroupRequest{
				Title:       &input.Title,
				Description: &input.Description,
			}
			resp, err := ts.Client.PostInventoriesIdPlanGroupsWithResponse(ctx, input.InventoryID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON201)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) AddPlanToGroupTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "add_plan_to_group",
			Description: "Add a plan to a plan group.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"group_id": {
						"type": "string",
						"description": "The ID of the plan group."
					},
					"plan_id": {
						"type": "string",
						"description": "The ID of the plan to add."
					}
				},
				"required": ["group_id", "plan_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				GroupID string `json:"group_id"`
				PlanID  string `json:"plan_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersAddPlanToGroupRequest{
				PlanId: &input.PlanID,
			}
			resp, err := ts.Client.PostPlanGroupsIdPlansWithResponse(ctx, input.GroupID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			return "success", nil
		},
	}
}

func (ts *ToolSet) CreateShoppingListFromPlanTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "create_shopping_list_from_plan",
			Description: "Create a shopping list from a plan.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"plan_id": {
						"type": "string",
						"description": "The ID of the plan."
					},
					"name": {
						"type": "string",
						"description": "The name of the new shopping list."
					}
				},
				"required": ["plan_id", "name"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				PlanID string `json:"plan_id"`
				Name   string `json:"name"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersCreateShoppingListRequest{
				Name: &input.Name,
			}
			resp, err := ts.Client.PostPlansIdShoppingListWithResponse(ctx, input.PlanID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON201)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) CreateShoppingListFromPlanGroupTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "create_shopping_list_from_plan_group",
			Description: "Create a shopping list from a plan group.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"group_id": {
						"type": "string",
						"description": "The ID of the plan group."
					},
					"name": {
						"type": "string",
						"description": "The name of the new shopping list."
					}
				},
				"required": ["group_id", "name"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				GroupID string `json:"group_id"`
				Name    string `json:"name"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersCreateShoppingListFromGroupRequest{
				Name: &input.Name,
			}
			resp, err := ts.Client.PostPlanGroupsIdShoppingListWithResponse(ctx, input.GroupID, body)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 201 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON201)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) ListPlansTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "list_plans",
			Description: "List all plans in the current inventory.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					}
				},
				"required": ["inventory_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			resp, err := ts.Client.GetInventoriesIdPlansWithResponse(ctx, input.InventoryID, nil)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON200)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) ListPlanGroupsTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "list_plan_groups",
			Description: "List all plan groups in the current inventory.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					}
				},
				"required": ["inventory_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			resp, err := ts.Client.GetInventoriesIdPlanGroupsWithResponse(ctx, input.InventoryID, nil)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON200)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) GetPlanTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "get_plan",
			Description: "Get details of a specific plan.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"plan_id": {
						"type": "string",
						"description": "The ID of the plan."
					}
				},
				"required": ["plan_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				PlanID string `json:"plan_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			resp, err := ts.Client.GetPlansIdWithResponse(ctx, input.PlanID)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON200)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) GetPlanGroupTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "get_plan_group",
			Description: "Get details of a specific plan group.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"group_id": {
						"type": "string",
						"description": "The ID of the plan group."
					}
				},
				"required": ["group_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				GroupID string `json:"group_id"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			resp, err := ts.Client.GetPlanGroupsIdWithResponse(ctx, input.GroupID)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed: status %d", resp.StatusCode())
			}
			b, _ := json.Marshal(resp.JSON200)
			return string(b), nil
		},
	}
}
