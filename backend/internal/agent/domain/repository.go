package domain

import "context"

type AgentRepository interface {
	Create(ctx context.Context, agent *Agent) error
	FindByID(ctx context.Context, id string) (*Agent, error)
	ListByType(ctx context.Context, agentType AgentType) ([]*Agent, error)
	Update(ctx context.Context, agent *Agent) error
	Delete(ctx context.Context, id string) error
}
