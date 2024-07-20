package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetAllArticlesResponse struct {
	Authors []Authors `json:"authors"`
}

// GetAllArticles godoc
//
// @Summary Получить все статьи
// @Description Возвращает все статьи, сгруппированные по авторам.
// @Tags Article
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} GetAllArticlesResponse
// @Failure 500 {object} ErrorResponse
// @Router /article/all [get]
func (s *Server) GetAllArticles(c *gin.Context) {
	articles, err := s.storage.GetAllArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	authorMap := make(map[uint32]*Authors)

	for _, article := range articles {
		if _, ok := authorMap[article.AuthorId]; !ok {
			authorMap[article.AuthorId] = &Authors{
				Username: article.Author.Username,
				Articles: []Article{},
			}
		}
		authorMap[article.AuthorId].Articles = append(authorMap[article.AuthorId].Articles, Article{
			Title:   article.Title,
			Content: article.Content,
		})
	}

	response := GetAllArticlesResponse{
		Authors: make([]Authors, 0, len(authorMap)),
	}

	for _, author := range authorMap {
		response.Authors = append(response.Authors, *author)
	}

	c.JSON(http.StatusOK, response)
}
