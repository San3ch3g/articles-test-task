package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAuthor godoc
//
// @Summary      Получение задач пользователя
// @Description  Получение задач пользователя по номеру паспорта
// @Tags         Author
// @Produce      json
// @Router       /author [get]
func (s *Server) GetAuthor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Получение авторов"})
}
