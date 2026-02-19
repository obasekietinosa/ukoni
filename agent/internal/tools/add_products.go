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
			Description: "Add one or more canonical products to the inventory. Checks for duplicates before adding.",
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
				// Search by name to check for existing product
				searchParams := &client.GetInventoriesIdCanonicalProductsParams{
					Search: &p.Name,
				}
				searchResp, err := ts.Client.GetInventoriesIdCanonicalProductsWithResponse(ctx, input.InventoryID, searchParams)
				if err != nil {
					results = append(results, fmt.Sprintf("Error searching for %s: %v", p.Name, err))
					continue
				}
				if searchResp.StatusCode() != 200 {
					results = append(results, fmt.Sprintf("Error searching for %s: status %d", p.Name, searchResp.StatusCode()))
					continue
				}

				exists := false
				if searchResp.JSON200 != nil {
					for _, existing := range *searchResp.JSON200 {
						if existing.Name != nil && strings.EqualFold(*existing.Name, p.Name) {
							exists = true
							break
						}
					}
				}

				if exists {
					results = append(results, fmt.Sprintf("Product '%s' already exists (skipped)", p.Name))
					continue
				}

				// Create product
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
