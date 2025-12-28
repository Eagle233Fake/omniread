package domain

type AgentType string

const (
	AgentTypeCharacter  AgentType = "character"
	AgentTypeReader     AgentType = "reader"
	AgentTypeHistorical AgentType = "historical"
)

// AgentConfig 定义 Agent 的动态配置
type AgentConfig struct {
	EnableInternet   bool     `bson:"enable_internet" json:"enable_internet"`
	KnowledgeBaseIDs []string `bson:"knowledge_base_ids" json:"knowledge_base_ids"`
	Model            string   `bson:"model" json:"model"` // 指定特定模型
}

// AgentProfile 定义 Agent 的形象设定
type AgentProfile struct {
	// 通用字段
	Avatar   string `bson:"avatar" json:"avatar"`
	Language string `bson:"language" json:"language"`

	// 书中人物特有
	BookName string `bson:"book_name,omitempty" json:"book_name,omitempty"`
	RoleName string `bson:"role_name,omitempty" json:"role_name,omitempty"`

	// 读者群体特有
	Profession string `bson:"profession,omitempty" json:"profession,omitempty"`
	Interest   string `bson:"interest,omitempty" json:"interest,omitempty"`

	// 历史名人特有
	HistoricalEra string `bson:"historical_era,omitempty" json:"historical_era,omitempty"`

	// 自定义/补充设定 (用户输入或爬取的文本)
	Bio          string `bson:"bio" json:"bio"`
	CustomPrompt string `bson:"custom_prompt,omitempty" json:"custom_prompt,omitempty"`
}
