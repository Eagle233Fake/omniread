package repo

import (
	"context"

	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var AgentRepoSet = wire.NewSet(NewAgentRepo)

type AgentRepo struct {
	coll *mongo.Collection
}

func NewAgentRepo(db *mongo.Database) domain.AgentRepository {
	return &AgentRepo{
		coll: db.Collection("agents"),
	}
}

func (r *AgentRepo) Create(ctx context.Context, agent *domain.Agent) error {
	if agent.ID == "" {
		agent.ID = primitive.NewObjectID().Hex()
	}
	_, err := r.coll.InsertOne(ctx, agent)
	return err
}

func (r *AgentRepo) FindByID(ctx context.Context, id string) (*domain.Agent, error) {
	var agent domain.Agent
	filter := bson.M{"_id": id}
	err := r.coll.FindOne(ctx, filter).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepo) ListByType(ctx context.Context, agentType domain.AgentType) ([]*domain.Agent, error) {
	filter := bson.M{"type": agentType}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []*domain.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *AgentRepo) Update(ctx context.Context, agent *domain.Agent) error {
	filter := bson.M{"_id": agent.ID}
	update := bson.M{
		"$set": agent,
	}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *AgentRepo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.coll.DeleteOne(ctx, filter)
	return err
}
