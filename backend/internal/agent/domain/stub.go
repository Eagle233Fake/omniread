package domain

import "context"

type StubAgentRepository struct{}

func (r *StubAgentRepository) Create(ctx context.Context, agent *Agent) error { return nil }
func (r *StubAgentRepository) FindByID(ctx context.Context, id string) (*Agent, error) {
	return &Agent{
		ID:   "stub-id",
		Name: "Stub Agent",
		Type: AgentTypeCharacter,
		Config: AgentConfig{
			EnableInternet: true,
		},
	}, nil
}
func (r *StubAgentRepository) ListByType(ctx context.Context, agentType AgentType) ([]*Agent, error) {
	return nil, nil
}
func (r *StubAgentRepository) Update(ctx context.Context, agent *Agent) error { return nil }
func (r *StubAgentRepository) Delete(ctx context.Context, id string) error { return nil }
