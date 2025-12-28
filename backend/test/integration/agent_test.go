package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Eagle233Fake/omniread/backend/api/handler/agent"
	"github.com/Eagle233Fake/omniread/backend/internal/agent/domain"
	"github.com/stretchr/testify/assert"
)

type CloseNotifyRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func (c *CloseNotifyRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func NewCloseNotifyRecorder() *CloseNotifyRecorder {
	return &CloseNotifyRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closed:           make(chan bool),
	}
}

func TestAgentChatFlow(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI")
	}

	r := setupRouter()

	// 1. Create Agent
	createReq := agent.CreateAgentRequest{
		Name:        "Test Historian",
		Type:        domain.AgentTypeHistorical,
		Description: "A test historian agent",
		Config:      domain.AgentConfig{EnableInternet: true},
		Profile:     domain.AgentProfile{Bio: "I am a test historian.", HistoricalEra: "Modern"},
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/v1/agents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createResp struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.Nil(t, err)
	assert.NotEmpty(t, createResp.Data.ID)
	agentID := createResp.Data.ID

	// 2. Chat with Agent
	// Use a simple message
	chatReq := agent.ChatRequest{
		AgentID: agentID,
		Message: "Hello, just say hi.",
	}
	body, _ = json.Marshal(chatReq)
	req, _ = http.NewRequest("POST", "/v1/agents/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w2 := NewCloseNotifyRecorder()
	r.ServeHTTP(w2, req)

	assert.Equal(t, http.StatusOK, w2.Code)
	// Validate SSE headers
	assert.Contains(t, w2.Header().Get("Content-Type"), "text/event-stream")

	// Validate body content (it's a stream)
	responseBody := w2.Body.String()
	// Depending on LLM response, it might be empty if LLM fails, but we expect at least headers or empty stream?
	// The handler writes "event: message"

	// Since LLM might be slow, the test might finish before stream?
	// No, ServeHTTP blocks until handler returns.
	// The handler blocks until stream.Recv() returns EOF or error.
	// So it should contain the full response.

	assert.Contains(t, responseBody, "event:message")
	assert.Contains(t, responseBody, "data:")
}
