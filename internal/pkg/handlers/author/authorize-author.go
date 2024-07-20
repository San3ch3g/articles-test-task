package author

import (
	"articleModule/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthorizeAuthor godoc
//
// @Summary Авторизация автора
// @Description Авторизация нового автора по имени и паролю
// @Tags Author
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Детали для авторизации автора"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /author/sign-in [post]
func (s *Server) AuthorizeAuthor(c *gin.Context) {
	var request AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	author, statusCode, err := s.storage.AuthorizeAuthor(request.Username, request.Password)
	if err != nil {
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}
	token, err := service.GenerateUserToken([]byte(s.cfg.Secret), author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
