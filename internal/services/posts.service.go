package services

import (
	"context"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"go.uber.org/zap/zapcore"
)

type IPostsService interface {
	GetPostsByUserID(ctx context.Context, req dto.GetPostsByUserIDRequest) (dto.GetPostsByUserIDResponse, error)
}

type postsService struct {
	postsRepo repo.IPostRepository
}

// GetPostsByUserID implements [IPostsService].
func (p *postsService) GetPostsByUserID(ctx context.Context, req dto.GetPostsByUserIDRequest) (dto.GetPostsByUserIDResponse, error) {
	posts, total, err := p.postsRepo.GetPostsByUserID(ctx, database.GetUserPostsWithCountParams{
		UserID: req.UserID,
		Limit:  int32(req.PageSize),
		Offset: int32(req.Page),
	})
	if err != nil {
		global.Logger.Error("failed to get posts by user id: ", zapcore.Field{Type: zapcore.ErrorType, Interface: err})
		return dto.GetPostsByUserIDResponse{}, err
	}

	res := dto.NewPagination(posts, total, req.Page, req.PageSize)

	return dto.GetPostsByUserIDResponse{
		Posts:     res.Items,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: res.TotalPage,
	}, nil
}

func NewPostsService(postRepo repo.IPostRepository) IPostsService {
	return &postsService{
		postsRepo: postRepo,
	}
}
