package insight

import (
	"context"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/errorx"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/repo"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var InsightServiceSet = wire.NewSet(NewInsightService)

type InsightService struct {
	sessionRepo  *repo.ReadingSessionRepo
	progressRepo *repo.ReadingProgressRepo
}

func NewInsightService(sessionRepo *repo.ReadingSessionRepo, progressRepo *repo.ReadingProgressRepo) *InsightService {
	return &InsightService{
		sessionRepo:  sessionRepo,
		progressRepo: progressRepo,
	}
}

func (s *InsightService) GetSummary(ctx context.Context, userID string) (*dto.InsightSummaryResp, error) {
	uid, _ := primitive.ObjectIDFromHex(userID)
	
	// Get stats for last 30 days
	end := time.Now()
	start := end.AddDate(0, 0, -30)
	
	sessions, err := s.sessionRepo.FindByUserAndDateRange(ctx, uid, start, end)
	if err != nil {
		return nil, errorx.New(500)
	}

	var totalDuration int64
	dailyStatsMap := make(map[string]int64)
	
	for _, sess := range sessions {
		totalDuration += sess.Duration
		dailyStatsMap[sess.Date] += sess.Duration
	}

	// Calculate streaks (simplified)
	// TODO: Proper streak calculation requires sorted unique dates
	streak := 0
	
	// Count finished books
	progressList, err := s.progressRepo.ListByUser(ctx, uid)
	if err != nil {
		return nil, errorx.New(500)
	}
	finishedCount := 0
	for _, p := range progressList {
		if p.Status == "finished" || p.Progress >= 100 {
			finishedCount++
		}
	}

	var dailyStats []dto.DailyStat
	for i := 0; i < 30; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		dailyStats = append(dailyStats, dto.DailyStat{
			Date:     date,
			Duration: dailyStatsMap[date],
		})
	}

	return &dto.InsightSummaryResp{
		TotalReadingTime: totalDuration,
		TotalBooksRead:   finishedCount,
		CurrentStreak:    streak, // Placeholder
		DailyStats:       dailyStats,
	}, nil
}
