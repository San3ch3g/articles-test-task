package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteArticleRequest struct {
	ArticleId uint32 `form:"id"`
}

// DeleteArticle godoc
//
// @Summary Удаление статьи
// @Description Удаляет статью по её идентификатору.
// @Tags Article
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param id query uint32 true "Идентификатор статьи"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /article/:id [delete]
func (s *Server) DeleteArticle(c *gin.Context) {
	var request DeleteArticleRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	statusCode, err := s.storage.DeleteArticle(request.ArticleId)
	if err != nil {
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
