package tools

import (
	"ukoni/agent/pkg/client"
)

type ToolSet struct {
	Client *client.ClientWithResponses
}

func NewToolSet(c *client.ClientWithResponses) *ToolSet {
	return &ToolSet{Client: c}
}

func (ts *ToolSet) GetTools() []Tool {
	return []Tool{
		ts.ListInventoryProductsTool(),
		ts.SearchCanonicalProductsTool(),
		ts.ListShoppingListsTool(),
		ts.CreateShoppingListTool(),
		ts.AddShoppingListItemTool(),
		ts.AddCanonicalProductsTool(),
		ts.CreatePlanTool(),
		ts.AddPlanItemTool(),
		ts.CreatePlanGroupTool(),
		ts.AddPlanToGroupTool(),
		ts.CreateShoppingListFromPlanTool(),
		ts.CreateShoppingListFromPlanGroupTool(),
		ts.ListPlansTool(),
		ts.ListPlanGroupsTool(),
		ts.GetPlanTool(),
		ts.GetPlanGroupTool(),
	}
}
