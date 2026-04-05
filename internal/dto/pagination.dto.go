package dto

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

// response
type PaginationResponse[T any] struct {
	Items     T     `json:"items"`
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	PageSize  int64 `json:"page_size"`
	TotalPage int64 `json:"total_page"`
}

func NewPagination[T any](items T, total int64, page, pageSize int64) PaginationResponse[T] {
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))
	return PaginationResponse[T]{
		Items:     items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}
}

// request
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type PaginationRequest struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

func (r *PaginationRequest) GetOffset() int64 {
	return (r.Page - 1) * r.PageSize
}

func (r *PaginationRequest) normalize() {
	if r.Page < 1 {
		r.Page = DefaultPage
	}
	if r.PageSize < 1 {
		r.PageSize = DefaultPageSize
	}
	if r.PageSize > MaxPageSize {
		r.PageSize = MaxPageSize
	}
}

func GetPaginationRequest(c *gin.Context) (*PaginationRequest, error) {
	var req PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, response.ErrInvalidParam
	}

	req.normalize()

	return &req, nil
}
