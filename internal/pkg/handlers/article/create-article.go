package article

import (
	"articleModule/internal/pkg/models"
	"articleModule/internal/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateArticle godoc
//
// @Summary Создание новой статьи
// @Description Создает новую статью с заголовком и содержимым, предоставленными в запросе. Требуется авторизация.
// @Tags Article
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param article body CreateArticleRequest true "Данные для создания статьи"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /article [post]
func (s *Server) CreateArticle(c *gin.Context) {
	log.Println("CreateArticle handler called")

	var request CreateArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if len(request.Title) < 3 || len(request.Title) > 100 {
		log.Printf("Invalid title length: %s", request.Title)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "title must be between 3 and 100 characters"})
		return
	}
	if !service.IsValidText(request.Title) || !service.IsValidText(request.Content) {
		log.Printf("Invalid characters in title or content: %s, %s", request.Title, request.Content)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "title and content must contain only letters"})
		return
	}

	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		log.Println("Authorization token missing")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization token isn't exist"})
		return
	}

	token := fmt.Sprintf("%v", authorization)
	claims, err := service.GetTokenClaimsFromJWT(token, []byte(s.cfg.Secret))
	if err != nil {
		log.Printf("Error getting token claims: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	article := models.Article{
		Title:    request.Title,
		Content:  request.Content,
		AuthorId: claims.AuthorId,
	}

	log.Printf("Creating article: %s", article.Title)
	if err := s.storage.CreateArticle(article); err != nil {
		log.Printf("Error creating article: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Article created successfully: %s", article.Title)
	c.JSON(http.StatusCreated, SuccessResponse{Success: true})
}
