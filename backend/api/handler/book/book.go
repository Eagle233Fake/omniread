package book

import (
	"strconv"

	"github.com/Eagle233Fake/omniread/backend/api/handler"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/provider"
	"github.com/gin-gonic/gin"
)

// UploadBook handles book file upload
// @Summary Upload Book
// @Tags Book
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Book file (pdf/epub)"
// @Param title formData string false "Book Title"
// @Param author formData string false "Book Author"
// @Success 200 {object} dto.BookResp
// @Router /books/upload [post]
func UploadBook(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}

	var req dto.UploadBookReq
	req.Title = c.PostForm("title")
	req.Author = c.PostForm("author")
	req.Description = c.PostForm("description")

	// Get user ID from context (set by auth middleware)
	userID := c.GetString("uid")

	resp, err := provider.Get().BookService.UploadBook(c.Request.Context(), userID, file, &req)
	handler.PostProcess(c, req, resp, err)
}

// ListBooks
// @Summary List Books
// @Tags Book
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} dto.BookResp
// @Router /books [get]
func ListBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := provider.Get().BookService.ListBooks(c.Request.Context(), page, limit)
	handler.PostProcess(c, nil, resp, err)
}

// UpdateBook
// @Summary Update Book Details
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param body body dto.UpdateBookReq true "Update Info"
// @Success 200 {object} dto.Resp
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("uid")
	var req dto.UpdateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.PostProcess(c, nil, nil, err)
		return
	}

	err := provider.Get().BookService.UpdateBook(c.Request.Context(), userID, id, &req)
	handler.PostProcess(c, req, dto.Success(), err)
}

// GetBook
// @Summary Get Book Details
// @Tags Book
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookResp
// @Router /books/{id} [get]
func GetBook(c *gin.Context) {
	id := c.Param("id")
	resp, err := provider.Get().BookService.GetBook(c.Request.Context(), id)
	handler.PostProcess(c, nil, resp, err)
}
