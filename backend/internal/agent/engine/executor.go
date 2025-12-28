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

	// 绑定工具到 ChatModel
	if len(agentTools) > 0 {
		toolInfos := make([]*schema.ToolInfo, 0, len(agentTools))
		for _, t := range agentTools {
			info, err := t.Info(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get tool info: %w", err)
			}
			toolInfos = append(toolInfos, info)
		}

		if err := chatModel.BindTools(toolInfos); err != nil {
			return nil, fmt.Errorf("failed to bind tools to chat model: %w", err)
		}
	}

	// 构造输入消息 (System + History + User)
	input := make([]*schema.Message, 0, len(history)+2)
	input = append(input, schema.SystemMessage(sysPrompt))
	input = append(input, history...)
	input = append(input, schema.UserMessage(userMessage))

	// 5. 执行逻辑 (ReAct Pattern)
	// 如果没有工具，直接流式返回
	if len(agentTools) == 0 {
		stream, err := chatModel.Stream(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to stream: %w", err)
		}
		return e.wrapStream(stream, sessionID, userMessage), nil
	}

	// 如果有工具，先尝试 Generate 探测意图
	resp, err := chatModel.Generate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate: %w", err)
	}

	// 检查是否有工具调用
	if len(resp.ToolCalls) > 0 {
		// 创建 ToolsNode
		toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: agentTools,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create tools node: %w", err)
		}

		// 执行工具
		toolOutputs, err := toolsNode.Invoke(ctx, resp)
		if err != nil {
			return nil, fmt.Errorf("failed to invoke tools: %w", err)
		}

		// 追加历史
		input = append(input, resp)
		input = append(input, toolOutputs...)

		// 第二轮流式调用
		stream, err := chatModel.Stream(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to stream final response: %w", err)
		}
		return e.wrapStream(stream, sessionID, userMessage), nil
	}

	// 无工具调用，直接返回 Generate 的结果
	return e.wrapStream(&memoryStream{content: resp.Content}, sessionID, userMessage), nil
}

func (e *AgentExecutor) wrapStream(inner Stream, sessionID, userMessage string) Stream {
	// 异步保存用户消息到历史 (System Prompt 不保存)
	go func() {
		bgCtx := context.Background()
		_ = e.historyManager.AddMessage(bgCtx, sessionID, schema.UserMessage(userMessage))
	}()

	return &historyStreamReader{
		reader:         inner,
		historyManager: e.historyManager,
		sessionID:      sessionID,
	}
}

type memoryStream struct {
	content string
	sent    bool
}

func (s *memoryStream) Recv() (*schema.Message, error) {
	if s.sent {
		return nil, io.EOF
	}
	s.sent = true
	return &schema.Message{
		Role:    "assistant",
		Content: s.content,
	}, nil
}

func (s *memoryStream) Close() {}

type historyStreamReader struct {
	reader         Stream
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
