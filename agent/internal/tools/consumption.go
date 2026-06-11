package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ukoni/agent/pkg/client"
)

type RecordConsumptionResult struct {
	Event             *client.ModelsConsumptionEvent `json:"event,omitempty"`
	InventoryAdjusted bool                           `json:"inventory_adjusted"`
	Transaction       *client.ModelsTransaction      `json:"transaction,omitempty"`
}

func (ts *ToolSet) RecordConsumptionTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "record_consumption",
			Description: "Record product consumption in an inventory. When a product_variant_id is provided, this tool also reduces inventory quantity by creating a negative transaction unless adjust_inventory is false. Use canonical_product_id only when a specific stocked variant is unknown.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					},
					"canonical_product_id": {
						"type": "string",
						"description": "The canonical product ID for the consumed item."
					},
					"product_variant_id": {
						"type": "string",
						"description": "The specific product variant ID. Provide this when available so inventory quantity can be reduced."
					},
					"quantity": {
						"type": "number",
						"description": "The amount consumed. Must be greater than zero."
					},
					"unit": {
						"type": "string",
						"description": "Unit for the consumed quantity, such as kg, L, pcs, or servings."
					},
					"note": {
						"type": "string",
						"description": "Optional note about the consumption event."
					},
					"consumed_at": {
						"type": "string",
						"description": "Optional RFC3339 timestamp. Defaults to the current time."
					},
					"adjust_inventory": {
						"type": "boolean",
						"description": "Whether to reduce inventory when product_variant_id is provided. Defaults to true."
					}
				},
				"required": ["inventory_id", "quantity"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID        string  `json:"inventory_id"`
				CanonicalProductID string  `json:"canonical_product_id"`
				ProductVariantID   string  `json:"product_variant_id"`
				Quantity           float32 `json:"quantity"`
				Unit               string  `json:"unit"`
				Note               string  `json:"note"`
				ConsumedAt         string  `json:"consumed_at"`
				AdjustInventory    *bool   `json:"adjust_inventory"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}
			if input.InventoryID == "" {
				return "", fmt.Errorf("inventory_id is required")
			}
			if input.CanonicalProductID == "" && input.ProductVariantID == "" {
				return "", fmt.Errorf("canonical_product_id or product_variant_id is required")
			}
			if input.Quantity <= 0 {
				return "", fmt.Errorf("quantity must be greater than zero")
			}

			consumedAt := input.ConsumedAt
			if consumedAt == "" {
				consumedAt = time.Now().UTC().Format(time.RFC3339)
			} else if _, err := time.Parse(time.RFC3339, consumedAt); err != nil {
				return "", fmt.Errorf("consumed_at must be RFC3339: %w", err)
			}

			source := "agent"
			consumptionBody := client.HandlersCreateConsumptionRequest{
				Quantity:   &input.Quantity,
				Source:     &source,
				ConsumedAt: &consumedAt,
			}
			if input.CanonicalProductID != "" {
				consumptionBody.CanonicalProductId = &input.CanonicalProductID
			}
			if input.ProductVariantID != "" {
				consumptionBody.ProductVariantId = &input.ProductVariantID
			}
			if input.Unit != "" {
				consumptionBody.Unit = &input.Unit
			}
			if input.Note != "" {
				consumptionBody.Note = &input.Note
			}

			consumptionResp, err := ts.Client.PostInventoriesIdConsumptionEventsWithResponse(ctx, input.InventoryID, consumptionBody)
			if err != nil {
				return "", err
			}
			if consumptionResp.StatusCode() != 201 {
				return "", fmt.Errorf("failed to record consumption: status %d", consumptionResp.StatusCode())
			}

			result := RecordConsumptionResult{
				Event: consumptionResp.JSON201,
			}

			shouldAdjust := input.AdjustInventory == nil || *input.AdjustInventory
			if shouldAdjust && input.ProductVariantID != "" {
				negativeQuantity := -input.Quantity
				transactionBody := client.HandlersCreateTransactionRequest{
					TransactionDate: &consumedAt,
					Items: &[]client.HandlersCreateTransactionItemRequest{
						{
							ProductVariantId: &input.ProductVariantID,
							Quantity:         &negativeQuantity,
						},
					},
				}

				transactionResp, err := ts.Client.PostInventoriesIdTransactionsWithResponse(ctx, input.InventoryID, transactionBody)
				if err != nil {
					return "", err
				}
				if transactionResp.StatusCode() != 201 {
					return "", fmt.Errorf("recorded consumption but failed to adjust inventory: status %d", transactionResp.StatusCode())
				}

				result.InventoryAdjusted = true
				result.Transaction = transactionResp.JSON201
			}

			b, _ := json.Marshal(result)
			return string(b), nil
		},
	}
}

func (ts *ToolSet) ListConsumptionEventsTool() Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        "list_consumption_events",
			Description: "List recent consumption events for an inventory.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inventory_id": {
						"type": "string",
						"description": "The ID of the inventory."
					},
					"limit": {
						"type": "integer",
						"description": "Maximum number of events to return. Defaults to 20."
					},
					"offset": {
						"type": "integer",
						"description": "Number of events to skip. Defaults to 0."
					}
				},
				"required": ["inventory_id"]
			}`),
		},
		Execute: func(ctx context.Context, args json.RawMessage) (string, error) {
			var input struct {
				InventoryID string `json:"inventory_id"`
				Limit       int    `json:"limit"`
				Offset      int    `json:"offset"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("invalid arguments: %w", err)
			}
			if input.InventoryID == "" {
				return "", fmt.Errorf("inventory_id is required")
			}
			if input.Limit < 0 {
				return "", fmt.Errorf("limit must be greater than or equal to zero")
			}
			if input.Offset < 0 {
				return "", fmt.Errorf("offset must be greater than or equal to zero")
			}

			params := &client.GetInventoriesIdConsumptionEventsParams{}
			if input.Limit > 0 {
				params.Limit = &input.Limit
			}
			if input.Offset > 0 {
				params.Offset = &input.Offset
			}

			resp, err := ts.Client.GetInventoriesIdConsumptionEventsWithResponse(ctx, input.InventoryID, params)
			if err != nil {
				return "", err
			}
			if resp.StatusCode() != 200 {
				return "", fmt.Errorf("failed to list consumption events: status %d", resp.StatusCode())
			}

			if resp.JSON200 == nil {
				return "[]", nil
			}
			b, _ := json.Marshal(resp.JSON200)
			return string(b), nil
		},
	}
}
