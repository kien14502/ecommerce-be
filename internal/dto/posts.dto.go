package dto

import "github.com/kien14502/ecommerce-be/internal/database"

// response
type GetPostsByUserIDResponse struct {
	Posts     []database.Post `json:"posts"`
	Total     int64           `json:"total"`
	Page      int64           `json:"page"`
	PageSize  int64           `json:"page_size"`
	TotalPage int64           `json:"total_page"`
}

type PostResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// request
type GetPostsByUserIDRequest struct {
	UserID   string `uri:"user_id"`
	Page     int64  `form:"page"`
	PageSize int64  `form:"page_size"`
}
