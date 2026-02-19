package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"ukoni/agent/pkg/client"
)

func (ts *ToolSet) AddCanonicalProductsTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "add_canonical_products",
			Description: "Add one or more canonical products to the inventory. Use this only for new products that do not already exist in the inventory context.",
			Parameters: json.RawMessage(`{
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
									"type": "string",
									"description": "The name of the product."
								},
								"description": {
									"type": "string",
									"description": "Optional description."
								},
								"category_id": {
									"type": "string",
									"description": "Optional category ID."
								}
							},
							"required": ["name"]
						}
					}
				},
				"required": ["inventory_id", "products"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
				Products    []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					CategoryID  string `json:"category_id"`
				} `json:"products"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}

			var results []string
			for _, p := range input.Products {
				// Create product directly, assuming LLM has filtered existing ones.
				body := client.HandlersCanonicalProductRequest{
					Name:        &p.Name,
					Description: &p.Description,
				}
				if p.CategoryID != "" {
					body.CategoryId = &p.CategoryID
				}

				createResp, err := ts.Client.PostInventoriesIdCanonicalProductsWithResponse(ctx, input.InventoryID, body)
				if err != nil {
					results = append(results, fmt.Sprintf("Error creating %s: %v", p.Name, err))
					continue
				}
				if createResp.StatusCode() != 201 {
					results = append(results, fmt.Sprintf("Error creating %s: status %d", p.Name, createResp.StatusCode()))
					continue
				}
				results = append(results, fmt.Sprintf("Product '%s' created", p.Name))
			}

			if len(results) == 0 {
				return "No products processed.", nil
			}
			return strings.Join(results, "\n"), nil
		},
	}
}
