package engine

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/engine/memory"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/engine/tools"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// Stream defines the interface for chat stream
type Stream interface {
	Recv() (*schema.Message, error)
	Close()
}

type AgentExecutor struct {
	promptBuilder  *PromptBuilder
	historyManager *memory.RedisHistoryManager
	config         *config.Config
}

func NewAgentExecutor(cfg *config.Config) *AgentExecutor {
	// 初始化 Redis 历史管理器
	hm := memory.NewRedisHistoryManager(
		cfg.Redis.Host,
		cfg.Redis.Pass,
		0, // 默认 DB
	)

	return &AgentExecutor{
		promptBuilder:  NewPromptBuilder(),
		historyManager: hm,
		config:         cfg,
	}
}

// ChatStream 返回一个 Reader 供上层流式读取
func (e *AgentExecutor) ChatStream(ctx context.Context, agent *domain.Agent, userMessage string) (Stream, error) {
	// 1. 获取历史消息
	// 假设 sessionID 为 agentID (简单起见，实际应为 用户ID:书籍ID)
	sessionID := agent.ID
	history, err := e.historyManager.GetHistory(ctx, sessionID)
	if err != nil {
		// 降级：如果 Redis 失败，仅记录错误，不阻断流程
		fmt.Printf("failed to get history: %v\n", err)
		history = []*schema.Message{}
	}

	// 2. 构建 System Prompt
	sysPrompt := e.promptBuilder.BuildSystemPrompt(agent)

	// 3. 准备 Tools
	var agentTools []tool.BaseTool
	if agent.Config.EnableInternet {
		bochaTool, err := tools.NewBochaSearchTool(ctx, &tools.BochaSearchConfig{
			APIKey: e.config.Bocha.APIKey,
		})
		if err == nil {
			agentTools = append(agentTools, bochaTool)
		} else {
			fmt.Printf("failed to init bocha tool: %v\n", err)
		}
	}

	// 4. 初始化 ChatModel
	modelName := "gpt-3.5-turbo"
	if e.config.Model.Model != "" {
		modelName = e.config.Model.Model
	}
	if agent.Config.Model != "" {
		modelName = agent.Config.Model
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: e.config.Model.BaseURL,
		Model:   modelName,
		APIKey:  e.config.Model.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init chat model: %w", err)
	}

	// 5. 定义 Graph (Chain)
	// 使用 ToolsNode 来支持工具调用
	// 流程： Input -> ToolsNode (LLM + Tools) -> Output
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()

	if len(agentTools) > 0 {
		// 如果有工具，使用 ToolsNode
		toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: agentTools,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create tools node: %w", err)
		}
		// 提取 ToolInfo 用于绑定到 ChatModel
		toolInfos := make([]*schema.ToolInfo, 0, len(agentTools))
		for _, t := range agentTools {
			info, err := t.Info(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get tool info: %w", err)
			}
			toolInfos = append(toolInfos, info)
		}

		// 正确绑定工具到 ChatModel
		if err := chatModel.BindTools(toolInfos); err != nil {
			return nil, fmt.Errorf("failed to bind tools to chat model: %w", err)
		}

		chain.AppendChatModel(chatModel)
		chain.AppendToolsNode(toolsNode)
	} else {
		// 无工具，直接连接 ChatModel
		chain.AppendChatModel(chatModel)
	}

	// 编译 Graph
	runner, err := chain.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compile chain: %w", err)
	}

	// 6. 构造输入消息 (System + History + User)
	input := make([]*schema.Message, 0, len(history)+2)
	input = append(input, schema.SystemMessage(sysPrompt))
	input = append(input, history...)
	input = append(input, schema.UserMessage(userMessage))

	// 7. 执行并获取流
	stream, err := runner.Stream(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to stream: %w", err)
	}

	// 8. 异步保存用户消息到历史 (System Prompt 不保存)
	go func() {
		bgCtx := context.Background()
		_ = e.historyManager.AddMessage(bgCtx, sessionID, schema.UserMessage(userMessage))
	}()

	return &historyStreamReader{
		reader:         stream,
		historyManager: e.historyManager,
		sessionID:      sessionID,
	}, nil
}

type historyStreamReader struct {
	reader         *schema.StreamReader[*schema.Message]
	historyManager *memory.RedisHistoryManager
	sessionID      string
	sb             strings.Builder
}

func (r *historyStreamReader) Recv() (*schema.Message, error) {
	msg, err := r.reader.Recv()
	if err == io.EOF {
		r.saveHistory()
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	r.sb.WriteString(msg.Content)
	return msg, nil
}

func (r *historyStreamReader) Close() {
	r.reader.Close()
}

func (r *historyStreamReader) saveHistory() {
	content := r.sb.String()
	if content == "" {
		return
	}

	msg := &schema.Message{
		Role:    "assistant",
		Content: content,
	}

	ctx := context.Background()
	_ = r.historyManager.AddMessage(ctx, r.sessionID, msg)
}
