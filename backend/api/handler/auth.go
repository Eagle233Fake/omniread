package handler

import (
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/provider"
	"github.com/Eagle233Fake/omniread/backend/types/errno"
	"github.com/gin-gonic/gin"
)

// Register handles user registration
// @Summary User Registration
// @Description Register a new user with username, password, and optional fields
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterReq true "Registration Request"
// @Success 200 {object} dto.Resp "Success"
// @Failure 400 {object} dto.Resp "Invalid parameters"
// @Failure 500 {object} dto.Resp "Internal server error"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		PostProcess(c, nil, nil, err)
		return
	}

	// Use service from provider
	err := provider.Get().AuthService.Register(c.Request.Context(), &req)

	if err != nil {
		PostProcess(c, req, nil, err)
		return
	}

	resp := dto.Success()
	PostProcess(c, req, resp, nil)
}

// Login handles user authentication
// @Summary User Login
// @Description Login with username/email/phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginReq true "Login Request"
// @Success 200 {object} dto.LoginResp "Success"
// @Failure 400 {object} dto.Resp "Invalid parameters"
// @Failure 401 {object} dto.Resp "Authentication failed"
// @Failure 500 {object} dto.Resp "Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		PostProcess(c, nil, nil, err)
		return
	}

	// Use service from provider
	data, err := provider.Get().AuthService.Login(c.Request.Context(), &req)

	if err != nil {
		PostProcess(c, req, nil, err)
		return
	}

	// Manual response construction not needed if we follow standard structure
	// But per user request to follow dto.go and not construct manually
	// We will use a structure that embeds Resp and the data
	resp := &struct {
		*dto.Resp
		*dto.LoginResp
	}{
		Resp:      dto.Success(),
		LoginResp: data,
	}

	// Override success code if needed, though Success() sets it to 0
	resp.Code = errno.SuccessCode

	PostProcess(c, req, resp, nil)
}
