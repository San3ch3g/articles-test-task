package author

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeleteAuthorRequest struct {
	Id uint32 `form:"id" binding:"required"`
}

// DeleteAuthor godoc
//
// @Summary Удаление автора
// @Description Удаление авторов по id
// @Tags Author
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param id query uint32 true "ID Автора"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /author/:id [delete]
func (s *Server) DeleteAuthor(c *gin.Context) {
	var request DeleteAuthorRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	statusCode, err := s.storage.DeleteAuthor(request.Id)
	if err != nil {
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
