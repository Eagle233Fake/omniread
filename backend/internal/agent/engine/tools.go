package engine

import "context"

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input string) (string, error)
}

// SearchTool 模拟联网搜索
type SearchTool struct{}

func (t *SearchTool) Name() string { return "internet_search" }
func (t *SearchTool) Description() string { return "Search the internet for real-time information." }
func (t *SearchTool) Execute(ctx context.Context, input string) (string, error) {
	return "Mock search result for: " + input, nil
}
