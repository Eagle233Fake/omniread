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

// GetProfile
// @Summary Get User Profile
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Router /user/profile [get]
func GetProfile(c *gin.Context) {
	userID := c.GetString("uid")
	user, err := provider.Get().AuthService.GetUser(c.Request.Context(), userID)
	PostProcess(c, nil, user, err)
}

// UpdateProfile
// @Summary Update User Profile
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileReq true "Update Request"
// @Success 200 {object} dto.Resp
// @Router /user/profile [put]
func UpdateProfile(c *gin.Context) {
	userID := c.GetString("uid")
	var req dto.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		PostProcess(c, nil, nil, err)
		return
	}
	err := provider.Get().AuthService.UpdateProfile(c.Request.Context(), userID, &req)
	PostProcess(c, req, dto.Success(), err)
}

// ChangePassword
// @Summary Change User Password
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordReq true "Change Password Request"
// @Success 200 {object} dto.Resp
// @Router /user/password [put]
func ChangePassword(c *gin.Context) {
	userID := c.GetString("uid")
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		PostProcess(c, nil, nil, err)
		return
	}
	err := provider.Get().AuthService.ChangePassword(c.Request.Context(), userID, &req)
	PostProcess(c, req, dto.Success(), err)
}

// UpdatePreferences
// @Summary Update User Preferences
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.UpdatePreferencesReq true "Update Preferences Request"
// @Success 200 {object} dto.Resp
// @Router /user/preferences [put]
func UpdatePreferences(c *gin.Context) {
	userID := c.GetString("uid")
	var req dto.UpdatePreferencesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		PostProcess(c, nil, nil, err)
		return
	}
	err := provider.Get().AuthService.UpdatePreferences(c.Request.Context(), userID, &req)
	PostProcess(c, req, dto.Success(), err)
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
