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
		ts.ListShoppingListsTool(),
		ts.CreateShoppingListTool(),
		ts.AddShoppingListItemTool(),
	}
}
