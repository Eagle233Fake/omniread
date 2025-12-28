package reading

import (
	"context"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/errorx"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/Eagle233Fake/omniread/backend/infra/repo"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ReadingServiceSet = wire.NewSet(NewReadingService)

type ReadingService struct {
	progressRepo *repo.ReadingProgressRepo
	sessionRepo  *repo.ReadingSessionRepo
}

func NewReadingService(progressRepo *repo.ReadingProgressRepo, sessionRepo *repo.ReadingSessionRepo) *ReadingService {
	return &ReadingService{
		progressRepo: progressRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s *ReadingService) UpdateProgress(ctx context.Context, userID string, req *dto.UpdateProgressReq) (*dto.ProgressResp, error) {
	uid, _ := primitive.ObjectIDFromHex(userID)
	bid, err := primitive.ObjectIDFromHex(req.BookID)
	if err != nil {
		return nil, errorx.New(400)
	}

	progress := &model.ReadingProgress{
		UserID:     uid,
		BookID:     bid,
		Progress:   req.Progress,
		CurrentLoc: req.CurrentLoc,
		Status:     req.Status,
		UpdatedAt:  time.Now(),
	}

	if err := s.progressRepo.Save(ctx, progress); err != nil {
		return nil, errorx.New(500)
	}

	return &dto.ProgressResp{
		BookID:     req.BookID,
		Progress:   req.Progress,
		CurrentLoc: req.CurrentLoc,
		Status:     req.Status,
		UpdatedAt:  progress.UpdatedAt,
	}, nil
}

func (s *ReadingService) GetProgress(ctx context.Context, userID, bookID string) (*dto.ProgressResp, error) {
	uid, _ := primitive.ObjectIDFromHex(userID)
	bid, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, errorx.New(400)
	}

	p, err := s.progressRepo.Find(ctx, uid, bid)
	if err != nil {
		return nil, errorx.New(500)
	}
	if p == nil {
		return nil, nil // No progress yet
	}

	return &dto.ProgressResp{
		BookID:     bookID,
		Progress:   p.Progress,
		CurrentLoc: p.CurrentLoc,
		Status:     p.Status,
		UpdatedAt:  p.UpdatedAt,
	}, nil
}

func (s *ReadingService) RecordSession(ctx context.Context, userID string, req *dto.ReadingSessionReq) error {
	uid, _ := primitive.ObjectIDFromHex(userID)
	bid, err := primitive.ObjectIDFromHex(req.BookID)
	if err != nil {
		return errorx.New(400)
	}

	session := &model.ReadingSession{
		UserID:    uid,
		BookID:    bid,
		StartTime: time.Unix(req.StartTime, 0),
		EndTime:   time.Unix(req.EndTime, 0),
		Duration:  req.Duration,
		Date:      time.Unix(req.StartTime, 0).Format("2006-01-02"),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return errorx.New(500)
	}
	return nil
}
