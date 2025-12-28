package reading

import (
	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/provider"
	"github.com/gin-gonic/gin"
)

// UpdateProgress
// @Summary Update Reading Progress
// @Tags Reading
// @Accept json
// @Produce json
// @Param request body dto.UpdateProgressReq true "Progress Data"
// @Success 200 {object} dto.ProgressResp
// @Router /reading/progress [post]
func UpdateProgress(c *gin.Context) {
	var req dto.UpdateProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}
	userID := c.GetString("uid")
	resp, err := provider.Get().ReadingService.UpdateProgress(c.Request.Context(), userID, &req)
	handler.PostProcess(c, req, resp, err)
}

// GetProgress
// @Summary Get Reading Progress
// @Tags Reading
// @Param book_id query string true "Book ID"
// @Success 200 {object} dto.ProgressResp
// @Router /reading/progress [get]
func GetProgress(c *gin.Context) {
	bookID := c.Query("book_id")
	userID := c.GetString("uid")
	resp, err := provider.Get().ReadingService.GetProgress(c.Request.Context(), userID, bookID)
	handler.PostProcess(c, nil, resp, err)
}

// RecordSession
// @Summary Record Reading Session
// @Tags Reading
// @Accept json
// @Produce json
// @Param request body dto.ReadingSessionReq true "Session Data"
// @Success 200 {object} dto.Resp
// @Router /reading/session [post]
func RecordSession(c *gin.Context) {
	var req dto.ReadingSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}
	userID := c.GetString("uid")
	err := provider.Get().ReadingService.RecordSession(c.Request.Context(), userID, &req)
	handler.PostProcess(c, req, dto.Success(), err)
}

// GetInsightSummary
// @Summary Get Reading Insight Summary
// @Tags Insight
// @Success 200 {object} dto.InsightSummaryResp
// @Router /insight/summary [get]
func GetInsightSummary(c *gin.Context) {
	userID := c.GetString("uid")
	resp, err := provider.Get().InsightService.GetSummary(c.Request.Context(), userID)
	handler.PostProcess(c, nil, resp, err)
}
