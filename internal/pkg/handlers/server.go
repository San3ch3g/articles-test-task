package handlers

import (
	"articleModule/internal/pkg/handlers/article"
	"articleModule/internal/pkg/handlers/author"
	"articleModule/internal/pkg/service"
	"articleModule/internal/pkg/storage/pg"
	"articleModule/internal/utils/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
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

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/author/sign-in" || c.Request.URL.Path == "/author/sign-up" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		_, err := service.ValidateToken(secret, authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Server) initRoutes() {
	authorHandlers := author.NewAuthorServer(s.storage, s.cfg)
	articleHandlers := article.NewArticleServer(s.storage, s.cfg)

	authGroup := s.router.Group("/", AuthMiddleware([]byte(s.cfg.Secret)))

	authorGroup := authGroup.Group("/author")
	{
		authorGroup.POST("/sign-in", authorHandlers.AuthorizeAuthor)
		authorGroup.POST("/sign-up", authorHandlers.RegisterAuthor)
		authorGroup.DELETE("/:id", authorHandlers.DeleteAuthor)
	}

	articleGroup := authGroup.Group("/article")
	{
		articleGroup.GET("/all", articleHandlers.GetAllArticles)
		articleGroup.POST("", articleHandlers.CreateArticle)
		articleGroup.DELETE("/:id", articleHandlers.DeleteArticle)
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
