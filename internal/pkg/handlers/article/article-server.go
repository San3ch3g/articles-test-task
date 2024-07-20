package article

import (
	"articleModule/internal/pkg/storage/pg"
	"articleModule/internal/utils/config"
)

type Server struct {
	storage *pg.Storage
	cfg     *config.Config
}

func NewArticleServer(storage *pg.Storage, cfg *config.Config) *Server {
	return &Server{
		storage: storage,
		cfg:     cfg,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type Authors struct {
	Username string    `json:"username"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
