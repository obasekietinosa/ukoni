package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (ts *ToolSet) ListInventoryProductsTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "list_inventory_products",
			Description: "List all products in the current inventory.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory to list products from."
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

			resp, err := ts.Client.GetInventoriesIdInventoryProductsWithResponse(ctx, input.InventoryID)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed to list products: status %d", resp.StatusCode())
			}

			if resp.JSON200 != nil {
                // Just dump it
				b, _ := json.Marshal(resp.JSON200)
				return string(b), nil
			}
            return "[]", nil
		},
	}
}
