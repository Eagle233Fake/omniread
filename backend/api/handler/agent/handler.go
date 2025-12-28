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

func (h *AgentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handler.PostProcess(c, nil, nil, errors.New("id is required"))
		return
	}

	var req UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}

	// First get the existing agent to preserve fields if needed, 
	// or we just overwrite. For simplicity, let's construct the object.
	// A better approach is to fetch, update fields, then save.
	// But since we are likely replacing the role configuration, overwriting 
	// the config/profile is exactly what we want.
	
	// We need to fetch the agent to ensure it exists and preserve created_at etc.
	// But service.UpdateAgent only updates what we give it? 
	// The repo implementation: $set: agent. This means it replaces fields provided in the struct.
	// But the struct has zero values for missing fields. 
	// Ideally we should fetch first.
	
	existing, err := h.service.GetAgent(c.Request.Context(), id)
	if err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}
	if existing == nil {
		handler.PostProcess(c, nil, nil, errors.New("agent not found"))
		return
	}

	// Update fields
	if req.Name != "" { existing.Name = req.Name }
	if req.Type != "" { existing.Type = req.Type }
	if req.Description != "" { existing.Description = req.Description }
	// For nested structs Config and Profile, we likely want to replace them if provided,
	// or merge them. Given the requirement is "switch role", replacing Profile is correct.
	// We assume if they are provided in JSON (even empty), they are meant to be updated.
	// However, Go's zero value makes it hard to distinguish "missing" vs "empty".
	// Let's assume the frontend sends the full profile for the new role.
	existing.Config = req.Config
	existing.Profile = req.Profile

	if err := h.service.UpdateAgent(c.Request.Context(), existing); err != nil {
		handler.PostProcess(c, req, nil, err)
		return
	}

	handler.PostProcess(c, req, nil, nil)
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
