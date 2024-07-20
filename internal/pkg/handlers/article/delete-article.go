package article

import (
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println("DeleteArticle handler called")

	var request DeleteArticleRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Printf("Error binding query parameters: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Attempting to delete article with ID: %d", request.ArticleId)
	statusCode, err := s.storage.DeleteArticle(request.ArticleId)
	if err != nil {
		log.Printf("Error deleting article with ID %d: %v", request.ArticleId, err)
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Article with ID %d deleted successfully", request.ArticleId)
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
