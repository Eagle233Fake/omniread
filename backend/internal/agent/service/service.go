package service

import (
	"context"
	"errors"
	"time"

	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/engine"
)

type AgentService struct {
	repo     domain.AgentRepository
	executor *engine.AgentExecutor
}

func NewAgentService(repo domain.AgentRepository, cfg *config.Config) *AgentService {
	return &AgentService{
		repo:     repo,
		executor: engine.NewAgentExecutor(cfg),
	}
}

func (s *AgentService) CreateAgent(ctx context.Context, agent *domain.Agent) error {
	if agent.Name == "" {
		return errors.New("agent name is required")
	}
	agent.CreatedAt = time.Now().Unix()
	agent.UpdatedAt = time.Now().Unix()
	return s.repo.Create(ctx, agent)
}

func (s *AgentService) GetAgent(ctx context.Context, id string) (*domain.Agent, error) {
	return s.repo.FindByID(ctx, id)
}

// ChatStream 建立流式对话
func (s *AgentService) ChatStream(ctx context.Context, agentID string, message string) (engine.Stream, error) {
	agent, err := s.repo.FindByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("agent not found")
	}

	return s.executor.ChatStream(ctx, agent, message)
}
