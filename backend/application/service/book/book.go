package book

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Boyuan-IT-Club/go-kit/errorx"
	"github.com/Eagle233Fake/omniread/backend/application/dto"
	"github.com/Eagle233Fake/omniread/backend/infra/model"
	"github.com/Eagle233Fake/omniread/backend/infra/oss"
	"github.com/Eagle233Fake/omniread/backend/infra/repo"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BookServiceSet = wire.NewSet(NewBookService)

type BookService struct {
	bookRepo *repo.BookRepo
	oss      *oss.OSSClient
}

func NewBookService(bookRepo *repo.BookRepo, oss *oss.OSSClient) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		oss:      oss,
	}
}

func (s *BookService) UploadBook(ctx context.Context, userID string, file *multipart.FileHeader, req *dto.UploadBookReq) (*dto.BookResp, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".epub" {
		return nil, errorx.New(400)
	}

	// Save file to OSS
	src, err := file.Open()
	if err != nil {
		return nil, errorx.New(500)
	}
	defer src.Close()

	objectName := uuid.New().String() + ext
	contentType := "application/octet-stream"
	if ext == ".pdf" {
		contentType = "application/pdf"
	} else if ext == ".epub" {
		contentType = "application/epub+zip"
	}

	fileURL, err := s.oss.UploadFile(ctx, objectName, src, file.Size, contentType)
	if err != nil {
		return nil, errorx.New(500)
	}

	// Parse UserID
	uid, _ := primitive.ObjectIDFromHex(userID)

	// Create Book Record
	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Format:      strings.TrimPrefix(ext, "."),
		Size:        file.Size,
		FileURL:     fileURL,
		UploadBy:    uid,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if book.Title == "" {
		book.Title = strings.TrimSuffix(file.Filename, ext)
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, errorx.New(500)
	}

	return &dto.BookResp{
		ID:          book.ID.Hex(),
		Title:       book.Title,
		Author:      book.Author,
		FileURL:     book.FileURL,
		Format:      book.Format,
		Size:        book.Size,
		Description: book.Description,
		CreatedAt:   book.CreatedAt,
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, userID, bookID string, req *dto.UpdateBookReq) error {
	oid, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return errorx.New(400)
	}

	book, err := s.bookRepo.FindByID(ctx, oid)
	if err != nil {
		return errorx.New(404)
	}

	// Permission Check: only owner can update (add admin check later if needed)
	// Allow update if UploadBy is zero (legacy/migration support)
	if !book.UploadBy.IsZero() && book.UploadBy.Hex() != userID {
		return errorx.New(403)
	}

	// If it was a legacy book, take ownership
	if book.UploadBy.IsZero() {
		uid, _ := primitive.ObjectIDFromHex(userID)
		book.UploadBy = uid
	}
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return errorx.New(500)
	}
	return nil
}

func (s *BookService) ListBooks(ctx context.Context, page, limit int) ([]*dto.BookResp, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := int64((page - 1) * limit)

	books, err := s.bookRepo.List(ctx, int64(limit), offset)
	if err != nil {
		return nil, errorx.New(500)
	}

	var resp []*dto.BookResp
	for _, b := range books {
		resp = append(resp, &dto.BookResp{
			ID:          b.ID.Hex(),
			Title:       b.Title,
			Author:      b.Author,
			FileURL:     b.FileURL,
			CoverURL:    b.CoverURL,
			Format:      b.Format,
			Size:        b.Size,
			Description: b.Description,
			CreatedAt:   b.CreatedAt,
		})
	}
	return resp, nil
}

func (s *BookService) GetBook(ctx context.Context, bookID string) (*dto.BookResp, error) {
	oid, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, errorx.New(400)
	}
	b, err := s.bookRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, errorx.New(404)
	}
	return &dto.BookResp{
		ID:          b.ID.Hex(),
		Title:       b.Title,
		Author:      b.Author,
		FileURL:     b.FileURL,
		CoverURL:    b.CoverURL,
		Format:      b.Format,
		Size:        b.Size,
		Description: b.Description,
		CreatedAt:   b.CreatedAt,
	}, nil
}
