package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

type PostController struct {
	postService services.IPostsService
}

func NewPostController(postService services.IPostsService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// GetPosts godoc
// @Summary      Get Posts By User ID
// @Description  Retrieve a paginated list of posts by user ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true  "User ID"
// @Param        page      query     int64   false  "Page number (default: 1)"
// @Param        page_size query     int64   false  "Number of items per page (default: 10, max: 100)"
// @Success      200  {object}  response.Response{data=dto.PaginationResponse[[]dto.PostResponse]} "Get posts successfully"
// @Failure      400  {object}  response.Response "Invalid parameters"
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      404  {object}  response.Response "User not found"
// @Failure      500  {object}  response.Response "Internal server error"
// @Security     BearerAuth
// @Router       /user/{user_id}/posts [get]
func (p *PostController) GetPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var in dto.GetPostsByUserIDRequest
	if err := c.ShouldBindUri(&in); err != nil {
		c.Error(response.ErrInvalidParam)
		return
	}
	pagination, err := dto.GetPaginationRequest(c)
	if err != nil {
		c.Error(response.ErrInvalidParam)
		return
	}
	res, err := p.postService.GetPostsByUserID(ctx, dto.GetPostsByUserIDRequest{
		UserID:   in.UserID,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Get posts successfully",
		Data:    res,
	})
}
