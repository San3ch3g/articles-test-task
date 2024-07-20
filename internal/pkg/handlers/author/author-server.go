package author

import (
	"articleModule/internal/pkg/storage/pg"
	"articleModule/internal/utils/config"
)

type Server struct {
	storage *pg.Storage
	cfg     *config.Config
}

func NewAuthorServer(storage *pg.Storage, cfg *config.Config) *Server {
	return &Server{
		storage: storage,
		cfg:     cfg,
	}
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
