package agent

import "github.com/Eagle233Fake/omniread/backend/internal/agent/domain"

type BaseResponse struct {
	Code int64  `json:"-"`
	Msg  string `json:"-"`
}

type CreateAgentRequest struct {
	Name        string              `json:"name" binding:"required"`
	Type        domain.AgentType    `json:"type" binding:"required"`
	Description string              `json:"description"`
	Config      domain.AgentConfig  `json:"config"`
	Profile     domain.AgentProfile `json:"profile"`
}

type CreateAgentResponse struct {
	BaseResponse
	ID string `json:"id"`
}

type UpdateAgentRequest struct {
	Name        string              `json:"name"`
	Type        domain.AgentType    `json:"type"`
	Description string              `json:"description"`
	Config      domain.AgentConfig  `json:"config"`
	Profile     domain.AgentProfile `json:"profile"`
}

type ChatRequest struct {
	AgentID string `json:"agent_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	BaseResponse
	Reply string `json:"reply"`
}
