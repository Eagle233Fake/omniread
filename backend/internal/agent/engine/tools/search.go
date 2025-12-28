package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// BochaSearchConfig 配置
type BochaSearchConfig struct {
	APIKey string
}

type bochaSearchTool struct {
	apiKey string
	client *http.Client
}

// NewBochaSearchTool 创建 Bocha 搜索工具
func NewBochaSearchTool(ctx context.Context, config *BochaSearchConfig) (tool.BaseTool, error) {
	return &bochaSearchTool{
		apiKey: config.APIKey,
		client: &http.Client{},
	}, nil
}

// BochaRequest Bocha API 请求结构
type BochaRequest struct {
	Query     string `json:"query"`
	Freshness string `json:"freshness,omitempty"` // oneDay, oneWeek, oneMonth, oneYear, noLimit
	Summary   bool   `json:"summary,omitempty"`
	Count     int    `json:"count,omitempty"`
}

// BochaResponse Bocha API 响应结构
type BochaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		WebPages struct {
			Value []struct {
				Name    string `json:"name"`
				URL     string `json:"url"`
				Summary string `json:"summary"`
				Snippet string `json:"snippet"`
			} `json:"value"`
		} `json:"webPages"`
	} `json:"data"`
}

// RunInput 定义工具输入的参数结构
type RunInput struct {
	Query string `json:"query" jsonschema:"description=The search query"`
}

// Info 返回工具元数据
func (t *bochaSearchTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	paramsOneOf, err := utils.GoStruct2ParamsOneOf[RunInput]()
	if err != nil {
		return nil, err
	}

	return &schema.ToolInfo{
		Name:        "internet_search",
		Desc:        "Search the internet for real-time information using Bocha Search. Use this tool when you need to find up-to-date information, news, or answers to questions about current events.",
		ParamsOneOf: paramsOneOf,
	}, nil
}

// InvokableRun 执行搜索
func (t *bochaSearchTool) InvokableRun(ctx context.Context, input string, opts ...tool.Option) (string, error) {
	// 解析输入
	var runInput RunInput
	if err := json.Unmarshal([]byte(input), &runInput); err != nil {
		// 尝试直接作为 query 处理（兼容简单字符串输入）
		runInput.Query = input
	}

	// 构造 Bocha 请求
	reqBody := BochaRequest{
		Query:     runInput.Query,
		Freshness: "noLimit",
		Summary:   true,
		Count:     10,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.bochaai.com/v1/web-search", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bocha api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var bochaResp BochaResponse
	if err := json.NewDecoder(resp.Body).Decode(&bochaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if bochaResp.Code != 200 && bochaResp.Code != 0 {
		return "", fmt.Errorf("bocha api error: %s", bochaResp.Msg)
	}

	// 格式化输出
	var sb strings.Builder
	for i, item := range bochaResp.Data.WebPages.Value {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Name))
		sb.WriteString(fmt.Sprintf("   URL: %s\n", item.URL))
		if item.Summary != "" {
			sb.WriteString(fmt.Sprintf("   Summary: %s\n", item.Summary))
		} else {
			sb.WriteString(fmt.Sprintf("   Snippet: %s\n", item.Snippet))
		}
		sb.WriteString("\n")
	}

	if sb.Len() == 0 {
		return "No results found.", nil
	}

	return sb.String(), nil
}
