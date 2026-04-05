package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

type IPostRepository interface {
	CreatePost(ctx context.Context, body database.CreatePostParams) error
	UpdatePost(ctx context.Context, body database.UpdatePostParams) error
	GetPostByID(ctx context.Context, postID string) (*database.Post, error)
	GetPostsByUserID(ctx context.Context, queries database.GetUserPostsWithCountParams) ([]database.Post, int64, error)
	DeletePostByID(ctx context.Context, body database.DeletePostParams) error
}

type postRepository struct {
	sqlc *database.Queries
}

// GetPostsByUserID implements [IPostRepository].
func (p *postRepository) GetPostsByUserID(ctx context.Context, queries database.GetUserPostsWithCountParams) ([]database.Post, int64, error) {
	rows, err := p.sqlc.GetUserPostsWithCount(ctx, queries)
	if err != nil {
		return []database.Post{}, 0, fmt.Errorf("get posts by user id failed: %w", err)
	}

	if len(rows) == 0 {
		return []database.Post{}, 0, nil
	}
	total := rows[0].TotalCount

	posts := make([]database.Post, len(rows))
	for i, row := range rows {
		posts[i] = database.Post{
			ID:         row.ID,
			UserID:     row.UserID,
			Content:    row.Content,
			Visibility: row.Visibility,
			CreatedAt:  row.CreatedAt,
		}
	}

	return posts, total, nil
}

// CreatePost implements [IPostRepository].
func (p *postRepository) CreatePost(ctx context.Context, body database.CreatePostParams) error {
	if err := p.sqlc.CreatePost(ctx, body); err != nil {
		return err
	}
	return nil
}

// DeletePostByID implements [IPostRepository].
func (p *postRepository) DeletePostByID(ctx context.Context, body database.DeletePostParams) error {
	if err := p.sqlc.DeletePost(ctx, body); err != nil {
		return err
	}
	return nil
}

// GetPostByID implements [IPostRepository].
func (p *postRepository) GetPostByID(ctx context.Context, postID string) (*database.Post, error) {
	post, err := p.sqlc.GetPost(ctx, postID)
	if err == sql.ErrNoRows {
		return nil, response.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get post by id failed: %w", err)
	}

	return &post, nil
}

// UpdatePost implements [IPostRepository].
func (p *postRepository) UpdatePost(ctx context.Context, body database.UpdatePostParams) error {

	if err := p.sqlc.UpdatePost(ctx, body); err != nil {
		return err
	}
	return nil
}

func NewPostRepository() IPostRepository {
	return &postRepository{
		sqlc: database.New(global.Mdbc),
	}
}
