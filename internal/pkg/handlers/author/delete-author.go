package author

import (
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println("DeleteAuthor handler called")

	var request DeleteAuthorRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		log.Printf("Error binding query parameters: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Attempting to delete author with ID: %d", request.Id)
	statusCode, err := s.storage.DeleteAuthor(request.Id)
	if err != nil {
		log.Printf("Error deleting author with ID %d: %v", request.Id, err)
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Author with ID %d deleted successfully", request.Id)
	c.JSON(http.StatusOK, SuccessResponse{Success: true})
}
