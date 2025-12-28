package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/stretchr/testify/assert"
)

func TestOmniReaderFlow(t *testing.T) {
	router := setupRouter()

	// 1. Register & Login
	token := registerAndLogin(t, router)

	// 2. Upload Book (Simulate PDF)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test_book.pdf")
	part.Write([]byte("%PDF-1.4 sample content"))
	writer.WriteField("title", "Integration Test Book")
	writer.WriteField("author", "Tester")
	writer.Close()

	req, _ := http.NewRequest("POST", "/books/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var bookResp struct {
		Code int          `json:"code"`
		Data dto.BookResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &bookResp)
	assert.NotEmpty(t, bookResp.Data.ID)
	assert.Equal(t, "Integration Test Book", bookResp.Data.Title)
	assert.Contains(t, bookResp.Data.FileURL, "http") // Should be OSS URL

	bookID := bookResp.Data.ID

	// 3. Update Reading Progress
	progressReq := dto.UpdateProgressReq{
		BookID:     bookID,
		Progress:   50.5,
		CurrentLoc: "page_10",
		Status:     "reading",
	}
	jsonBody, _ := json.Marshal(progressReq)
	req, _ = http.NewRequest("POST", "/reading/progress", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 4. Get Reading Progress
	req, _ = http.NewRequest("GET", fmt.Sprintf("/reading/progress?book_id=%s", bookID), nil)
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	
	var progressResp struct {
		Code int              `json:"code"`
		Data dto.ProgressResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &progressResp)
	assert.Equal(t, 50.5, progressResp.Data.Progress)
	assert.Equal(t, "page_10", progressResp.Data.CurrentLoc)
}

func TestOmniInsightFlow(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(t, router)

	// Create a dummy book first to get ID
	// Simplified upload reuse or mock ID if repo allows, but better to upload
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "insight_book.epub")
	part.Write([]byte("PK... sample epub"))
	writer.Close()
	req, _ := http.NewRequest("POST", "/books/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var bookResp struct {
		Data dto.BookResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &bookResp)
	bookID := bookResp.Data.ID

	// 1. Record Reading Session
	sessionReq := dto.ReadingSessionReq{
		BookID:    bookID,
		StartTime: time.Now().Unix() - 3600,
		EndTime:   time.Now().Unix(),
		Duration:  3600,
	}
	jsonBody, _ := json.Marshal(sessionReq)
	req, _ = http.NewRequest("POST", "/reading/session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 2. Get Insight Summary
	req, _ = http.NewRequest("GET", "/insight/summary", nil)
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	var summaryResp struct {
		Code int                    `json:"code"`
		Data dto.InsightSummaryResp `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &summaryResp)
	
	assert.GreaterOrEqual(t, summaryResp.Data.TotalReadingTime, int64(3600))
	assert.NotEmpty(t, summaryResp.Data.DailyStats)
}
