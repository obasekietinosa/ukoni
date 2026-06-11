package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ukoni/agent/pkg/client"
)

func TestConsumptionTools_Definition(t *testing.T) {
	ts := NewToolSet(&client.ClientWithResponses{})

	tests := []struct {
		name     string
		toolFunc func() Tool
		toolName string
	}{
		{"RecordConsumptionTool", ts.RecordConsumptionTool, "record_consumption"},
		{"ListConsumptionEventsTool", ts.ListConsumptionEventsTool, "list_consumption_events"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := tt.toolFunc()
			if tool.Definition.Name != tt.toolName {
				t.Errorf("expected tool name '%s', got '%s'", tt.toolName, tool.Definition.Name)
			}

			var schema map[string]interface{}
			if err := json.Unmarshal(tool.Definition.Parameters, &schema); err != nil {
				t.Errorf("invalid parameters schema for tool '%s': %v", tt.toolName, err)
			}
		})
	}
}

func TestRecordConsumptionTool_AdjustsInventoryForVariant(t *testing.T) {
	var sawConsumption bool
	var sawTransaction bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/inventories/inv-1/consumption-events":
			sawConsumption = true
			var body client.HandlersCreateConsumptionRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode consumption request: %v", err)
			}
			if body.ProductVariantId == nil || *body.ProductVariantId != "variant-1" {
				t.Fatalf("expected product_variant_id variant-1, got %#v", body.ProductVariantId)
			}
			if body.Quantity == nil || *body.Quantity != 2 {
				t.Fatalf("expected quantity 2, got %#v", body.Quantity)
			}
			if body.Source == nil || *body.Source != "agent" {
				t.Fatalf("expected source agent, got %#v", body.Source)
			}
			if body.ConsumedAt == nil || *body.ConsumedAt != "2026-06-10T12:00:00Z" {
				t.Fatalf("expected supplied consumed_at, got %#v", body.ConsumedAt)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"consumption-1","inventory_id":"inv-1","product_variant_id":"variant-1","quantity":2,"source":"agent","consumed_at":"2026-06-10T12:00:00Z"}`))

		case r.Method == http.MethodPost && r.URL.Path == "/inventories/inv-1/transactions":
			sawTransaction = true
			var body client.HandlersCreateTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode transaction request: %v", err)
			}
			if body.TransactionDate == nil || *body.TransactionDate != "2026-06-10T12:00:00Z" {
				t.Fatalf("expected transaction date to match consumed_at, got %#v", body.TransactionDate)
			}
			if body.Items == nil || len(*body.Items) != 1 {
				t.Fatalf("expected one transaction item, got %#v", body.Items)
			}
			item := (*body.Items)[0]
			if item.ProductVariantId == nil || *item.ProductVariantId != "variant-1" {
				t.Fatalf("expected transaction item variant-1, got %#v", item.ProductVariantId)
			}
			if item.Quantity == nil || *item.Quantity != -2 {
				t.Fatalf("expected negative transaction quantity -2, got %#v", item.Quantity)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"transaction-1","inventory_id":"inv-1","transaction_date":"2026-06-10T12:00:00Z"}`))

		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	c, err := client.NewClientWithResponses(server.URL)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tool := NewToolSet(c).RecordConsumptionTool()
	resultJSON, err := tool.Execute(context.Background(), json.RawMessage(`{
		"inventory_id": "inv-1",
		"product_variant_id": "variant-1",
		"quantity": 2,
		"unit": "pcs",
		"consumed_at": "2026-06-10T12:00:00Z"
	}`))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !sawConsumption {
		t.Fatal("expected consumption endpoint to be called")
	}
	if !sawTransaction {
		t.Fatal("expected transaction endpoint to be called")
	}

	var result RecordConsumptionResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if !result.InventoryAdjusted {
		t.Fatal("expected inventory_adjusted to be true")
	}
	if result.Event == nil || result.Transaction == nil {
		t.Fatalf("expected event and transaction in result, got %#v", result)
	}
}

func TestRecordConsumptionTool_DoesNotAdjustInventoryWithoutVariant(t *testing.T) {
	var sawTransaction bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost && r.URL.Path == "/inventories/inv-1/transactions" {
			sawTransaction = true
			t.Fatal("transaction endpoint should not be called for canonical-only consumption")
		}
		if r.Method != http.MethodPost || r.URL.Path != "/inventories/inv-1/consumption-events" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		var body client.HandlersCreateConsumptionRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode consumption request: %v", err)
		}
		if body.CanonicalProductId == nil || *body.CanonicalProductId != "canonical-1" {
			t.Fatalf("expected canonical_product_id canonical-1, got %#v", body.CanonicalProductId)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"consumption-1","inventory_id":"inv-1","canonical_product_id":"canonical-1","quantity":1,"source":"agent","consumed_at":"2026-06-10T12:00:00Z"}`))
	}))
	defer server.Close()

	c, err := client.NewClientWithResponses(server.URL)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tool := NewToolSet(c).RecordConsumptionTool()
	resultJSON, err := tool.Execute(context.Background(), json.RawMessage(`{
		"inventory_id": "inv-1",
		"canonical_product_id": "canonical-1",
		"quantity": 1,
		"consumed_at": "2026-06-10T12:00:00Z"
	}`))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if sawTransaction {
		t.Fatal("transaction endpoint should not have been called")
	}

	var result RecordConsumptionResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.InventoryAdjusted {
		t.Fatal("expected inventory_adjusted to be false")
	}
	if result.Event == nil {
		t.Fatal("expected event in result")
	}
}

func TestListConsumptionEventsTool_UsesPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet || r.URL.Path != "/inventories/inv-1/consumption-events" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "5" || r.URL.Query().Get("offset") != "10" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[{"id":"consumption-1","inventory_id":"inv-1","quantity":1,"source":"agent","consumed_at":"2026-06-10T12:00:00Z"}]`))
	}))
	defer server.Close()

	c, err := client.NewClientWithResponses(server.URL)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tool := NewToolSet(c).ListConsumptionEventsTool()
	resultJSON, err := tool.Execute(context.Background(), json.RawMessage(`{
		"inventory_id": "inv-1",
		"limit": 5,
		"offset": 10
	}`))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	var events []client.ModelsConsumptionEventDetail
	if err := json.Unmarshal([]byte(resultJSON), &events); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one event, got %d", len(events))
	}
}
