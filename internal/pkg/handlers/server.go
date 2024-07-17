package handlers

import (
	_ "articleModule/docs"
	"articleModule/internal/pkg/storage/pg"
	"articleModule/internal/utils/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router  *gin.Engine
	storage *pg.Storage
	cfg     *config.Config
}

func NewServer(storage *pg.Storage, cfg *config.Config) *Server {
	router := gin.Default()
	server := &Server{
		router:  router,
		storage: storage,
		cfg:     cfg,
	}
	server.initRoutes()
	return server
}

func (s *Server) InitSwagger() {
	swagURL := "http://localhost:8080/swagger/doc.json"
	url := ginSwagger.URL(swagURL)
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func (s *Server) initRoutes() {
	author := s.router.Group("/author")
	{
		author.GET("", s.GetAuthor)

	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
