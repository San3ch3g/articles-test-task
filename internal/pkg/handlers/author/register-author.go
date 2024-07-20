package author

import (
	"articleModule/internal/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterAuthor godoc
//
// @Summary Регистрация нового автора
// @Description Регистрация нового автора по имени и паролю
// @Tags Author
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Детали для регистрации автора"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /author/sign-up [post]
func (s *Server) RegisterAuthor(c *gin.Context) {
	var request AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if !service.IsValidText(request.Username) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprint("username must contain only Latin letters")})
		return
	}

	author, statusCode, err := s.storage.RegisterAuthor(request.Username, request.Password)
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
