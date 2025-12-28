package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to register and login a user for testing
func registerAndLogin(t *testing.T, router *gin.Engine) string {
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
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 2. Login
	loginReq := dto.LoginReq{
		Username: username,
		Password: "Password123",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp struct {
		Code int           `json:"code"`
		Data dto.LoginResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Data.Token
}
