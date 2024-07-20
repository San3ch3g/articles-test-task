package author

import (
	"articleModule/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println("RegisterAuthor handler called")

	var request AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if !service.IsValidText(request.Username) {
		log.Printf("Invalid username: %s", request.Username)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "username must contain only Latin letters"})
		return
	}

	author, statusCode, err := s.storage.RegisterAuthor(request.Username, request.Password)
	if err != nil {
		log.Printf("Error registering author: %v", err)
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	token, err := service.GenerateUserToken([]byte(s.cfg.Secret), author)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Author %s registered successfully", request.Username)
	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
