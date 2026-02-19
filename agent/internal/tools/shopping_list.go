package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"ukoni/agent/pkg/client"
)

func (ts *ToolSet) ListShoppingListsTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "list_shopping_lists",
			Description: "List all shopping lists in the current inventory.",
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

			resp, err := ts.Client.GetInventoriesIdShoppingListsWithResponse(ctx, input.InventoryID)
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

func (ts *ToolSet) CreateShoppingListTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "create_shopping_list",
			Description: "Create a new shopping list.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					},
					"name": {
						"type": "string",
						"description": "The name of the shopping list."
					},
                    "notes": {
                        "type": "string",
                        "description": "Optional notes."
                    }
				},
				"required": ["inventory_id", "name"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
				Name        string `json:"name"`
				Notes       string `json:"notes"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersCreateListRequest{
				Name:  &input.Name,
				Notes: &input.Notes,
			}
			resp, err := ts.Client.PostInventoriesIdShoppingListsWithResponse(ctx, input.InventoryID, body)
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

func (ts *ToolSet) AddShoppingListItemTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "add_shopping_list_item",
			Description: "Add an item to a shopping list.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"shopping_list_id": {
						"type": "string",
						"description": "The ID of the shopping list."
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
                    "notes": {
                        "type": "string",
                        "description": "Notes."
                    }
				},
				"required": ["shopping_list_id", "target_type", "target_id", "quantity"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				ShoppingListID string  `json:"shopping_list_id"`
				TargetType     string  `json:"target_type"`
				TargetID       string  `json:"target_id"`
				Quantity       float32 `json:"quantity"`
				Unit           string  `json:"unit"`
				Notes          string  `json:"notes"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			body := client.HandlersAddItemRequest{
				TargetType: &input.TargetType,
				TargetId:   &input.TargetID,
				Quantity:   &input.Quantity,
				Unit:       &input.Unit,
				Notes:      &input.Notes,
			}
			resp, err := ts.Client.PostShoppingListsIdItemsWithResponse(ctx, input.ShoppingListID, body)
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
