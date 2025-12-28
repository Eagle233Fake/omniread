package domain

// Agent 聚合根
type Agent struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	Name        string       `bson:"name" json:"name"`
	Type        AgentType    `bson:"type" json:"type"`
	Description string       `bson:"description" json:"description"`
	Config      AgentConfig  `bson:"config" json:"config"`
	Profile     AgentProfile `bson:"profile" json:"profile"`
	CreatedAt   int64        `bson:"created_at" json:"created_at"`
	UpdatedAt   int64        `bson:"updated_at" json:"updated_at"`
}
