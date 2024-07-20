package author

import (
	"articleModule/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println("AuthorizeAuthor handler called")

	var request AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Attempting to authorize author: %s", request.Username)
	author, statusCode, err := s.storage.AuthorizeAuthor(request.Username, request.Password)
	if err != nil {
		log.Printf("Error authorizing author: %v", err)
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Generating token for author: %s", request.Username)
	token, err := service.GenerateUserToken([]byte(s.cfg.Secret), author)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Author %s authorized successfully", request.Username)
	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
