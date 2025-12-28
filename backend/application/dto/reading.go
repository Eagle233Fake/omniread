package dto

import "time"

type UpdateProgressReq struct {
	BookID     string  `json:"book_id" binding:"required"`
	Progress   float64 `json:"progress"` // percentage
	CurrentLoc string  `json:"current_loc"`
	Status     string  `json:"status"` // reading, finished
}

type ProgressResp struct {
	BookID     string    `json:"book_id"`
	Progress   float64   `json:"progress"`
	CurrentLoc string    `json:"current_loc"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ReadingSessionReq struct {
	BookID    string `json:"book_id" binding:"required"`
	StartTime int64  `json:"start_time"` // unix timestamp
	EndTime   int64  `json:"end_time"`   // unix timestamp
	Duration  int64  `json:"duration"`   // seconds
}

// Insight DTOs
type DailyStat struct {
	Date     string `json:"date"`
	Duration int64  `json:"duration"` // seconds
}

type InsightSummaryResp struct {
	TotalReadingTime int64       `json:"total_reading_time"`
	TotalBooksRead   int         `json:"total_books_read"`
	CurrentStreak    int         `json:"current_streak"`
	DailyStats       []DailyStat `json:"daily_stats"`
}
