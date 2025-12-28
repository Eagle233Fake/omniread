package dto

import "time"

type UploadBookReq struct {
	Title       string `form:"title"`
	Author      string `form:"author"`
	Description string `form:"description"`
}

type UpdateBookReq struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

type BookResp struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	CoverURL    string    `json:"cover_url"`
	FileURL     string    `json:"file_url"`
	Format      string    `json:"format"`
	Size        int64     `json:"size"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
