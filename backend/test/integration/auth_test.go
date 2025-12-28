package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Eagle233Fake/omniread/backend/api/router"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/config"
	"github.com/Eagle233Fake/omniread/backend/provider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func setupRouter() *gin.Engine {
	// Setup config for testing
	os.Setenv("CONFIG_PATH", "../../etc/config_test.yaml")
	if _, err := config.NewConfig(); err != nil {
		panic(err)
	}

	provider.Init()
	return router.SetupRoutes()
}

func TestAuthFlow(t *testing.T) {
	// Skip integration test if MongoDB is not available or config is missing
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI")
	}

	r := setupRouter()

	// Unique username for each run
	username := "testuser_" + time.Now().Format("20060102150405")
	email := username + "@example.com"
	phone := "138" + time.Now().Format("01021504")

	// 1. Register
	registerReq := dto.RegisterReq{
		Username:  username,
		Password:  "Password123",
		Email:     email,
		Phone:     phone,
		Gender:    "male",
		Birthdate: "1990-01-01",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Login with username
	loginReq := dto.LoginReq{
		Username: username,
		Password: "Password123",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data dto.LoginResp `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Data.Token)
	assert.Equal(t, username, resp.Data.User.Username)

	// 3. Login with email
	loginReq.Username = email
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Verify Redis Cache
	cfg := config.GetConfig()
	rds := redis.MustNewRedis(*cfg.Redis)
	key := fmt.Sprintf("auth:token:%s", resp.Data.Token)
	val, err := rds.GetCtx(context.Background(), key)
	assert.Nil(t, err)
	assert.Equal(t, resp.Data.User.ID.Hex(), val)
}
