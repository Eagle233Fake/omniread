package engine

import (
	"fmt"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"strings"
)

// PromptBuilder 负责构建系统提示词
type PromptBuilder struct{}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (b *PromptBuilder) BuildSystemPrompt(agent *domain.Agent) string {
	var sb strings.Builder

	// 基础设定
	sb.WriteString(fmt.Sprintf("You are %s. ", agent.Name))
	if agent.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n", agent.Description))
	}

	// 根据类型追加特定指令
	switch agent.Type {
	case domain.AgentTypeCharacter:
		b.buildCharacterPrompt(&sb, agent)
	case domain.AgentTypeReader:
		b.buildReaderPrompt(&sb, agent)
	case domain.AgentTypeHistorical:
		b.buildHistoricalPrompt(&sb, agent)
	}

	// 追加用户自定义 Prompt
	if agent.Profile.CustomPrompt != "" {
		sb.WriteString("\n" + agent.Profile.CustomPrompt)
	}

	return sb.String()
}

func (b *PromptBuilder) buildCharacterPrompt(sb *strings.Builder, agent *domain.Agent) {
	sb.WriteString(fmt.Sprintf("You are a character from the book '%s'. ", agent.Profile.BookName))
	if agent.Profile.Bio != "" {
		sb.WriteString(fmt.Sprintf("Here is your background: %s\n", agent.Profile.Bio))
	}
	sb.WriteString("Stay in character at all times. Use the tone and vocabulary appropriate for your setting.")
}

func (b *PromptBuilder) buildReaderPrompt(sb *strings.Builder, agent *domain.Agent) {
	sb.WriteString(fmt.Sprintf("You are a %s interested in %s. ", agent.Profile.Profession, agent.Profile.Interest))
	sb.WriteString("Analyze the text from your professional perspective. Provide critical insights.")
}

func (b *PromptBuilder) buildHistoricalPrompt(sb *strings.Builder, agent *domain.Agent) {
	sb.WriteString(fmt.Sprintf("You are a historical figure from the %s era. ", agent.Profile.HistoricalEra))
	sb.WriteString("Speak in a manner consistent with your time period. Do not reference modern technology unless explicitly asked to compare.")
}
