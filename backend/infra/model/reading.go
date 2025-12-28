package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	CoverURL    string             `bson:"cover_url" json:"cover_url"`
	FileURL     string             `bson:"file_url" json:"file_url"`
	Format      string             `bson:"format" json:"format"` // pdf, epub
	Size        int64              `bson:"size" json:"size"`     // in bytes
	TotalPages  int                `bson:"total_pages" json:"total_pages"`
	Publisher   string             `bson:"publisher" json:"publisher"`
	Description string             `bson:"description" json:"description"`
	UploadBy    primitive.ObjectID `bson:"upload_by" json:"upload_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReadingProgress struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	BookID     primitive.ObjectID `bson:"book_id" json:"book_id"`
	Progress   float64            `bson:"progress" json:"progress"`     // 0-100 percentage
	CurrentLoc string             `bson:"current_loc" json:"current_loc"` // cfi for epub, page number for pdf
	Status     string             `bson:"status" json:"status"`         // reading, finished
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReadingSession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	BookID    primitive.ObjectID `bson:"book_id" json:"book_id"`
	StartTime time.Time          `bson:"start_time" json:"start_time"`
	EndTime   time.Time          `bson:"end_time" json:"end_time"`
	Duration  int64              `bson:"duration" json:"duration"` // in seconds
	Date      string             `bson:"date" json:"date"`         // YYYY-MM-DD for easier aggregation
}

type UserPreference struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	FontFamily string             `bson:"font_family" json:"font_family"`
	FontSize   int                `bson:"font_size" json:"font_size"`
	Theme      string             `bson:"theme" json:"theme"` // light, dark, sepia
	LineHeight float64            `bson:"line_height" json:"line_height"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
