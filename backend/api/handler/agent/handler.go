package agent

import (
	"errors"
	"io"

	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/service"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	service *service.AgentService
}

func NewAgentHandler(svc *service.AgentService) *AgentHandler {
	return &AgentHandler{service: svc}
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}

	agent := &domain.Agent{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Config:      req.Config,
		Profile:     req.Profile,
	}

	if err := h.service.CreateAgent(c.Request.Context(), agent); err != nil {
		handler.PostProcess(c, req, nil, err)
		return
	}

	resp := &CreateAgentResponse{
		ID: agent.ID,
	}
	handler.PostProcess(c, req, resp, nil)
}

func (h *AgentHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}

	// 获取流式 Reader
	stream, err := h.service.ChatStream(c.Request.Context(), req.AgentID, req.Message)
	if err != nil {
		handler.PostProcess(c, req, nil, err)
		return
	}
	defer stream.Close()

	// 设置 SSE 头部
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return false
		}
		if err != nil {
			// 流式过程中出错，发送错误事件并断开
			c.SSEvent("error", err.Error())
			return false
		}

		// 发送增量内容
		c.SSEvent("message", msg.Content)
		return true
	})
}
