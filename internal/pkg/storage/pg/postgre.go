package pg

import (
	"articleModule/internal/pkg/models"
	"articleModule/internal/utils/config"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Storage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func buildDSN(cfg *config.Config) string {
	//dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	return dsn
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return nil, err
	}
	log.Println("Database connection established")
	return db, nil
}

func MustNewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Author{})
	if err != nil {
		log.Println("table 'authors' already exists")
	} else {
		log.Println("table 'authors' migrated successfully")
	}
	err = db.AutoMigrate(&models.Article{})
	if err != nil {
		log.Println("table 'articles' already exists")
	} else {
		log.Println("table 'articles' migrated successfully")
	}
	return db
}

func (s *Storage) RegisterAuthor(username, password string) (models.Author, int, error) {
	log.Printf("Registering author: %s", username)

	var foundAuthor models.Author
	if err := s.db.Where("username = ?", username).First(&foundAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newAuthor models.Author
			newAuthor.Username = username
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Error hashing password: %v", err)
				return models.Author{}, http.StatusInternalServerError, err
			}
			newAuthor.Password = string(hashedPassword)
			if err := s.db.Create(&newAuthor).Error; err != nil {
				log.Printf("Error creating author: %v", err)
				return models.Author{}, http.StatusInternalServerError, err
			}
			log.Printf("Author %s registered successfully", username)
			return newAuthor, http.StatusCreated, nil
		} else {
			log.Printf("Error finding author: %v", err)
			return models.Author{}, http.StatusConflict, err
		}
	}
	log.Printf("Author %s already exists", username)
	return models.Author{}, http.StatusConflict, fmt.Errorf("author with username [%s] already exists", username)
}

func (s *Storage) AuthorizeAuthor(username, password string) (models.Author, int, error) {
	log.Printf("Authorizing author: %s", username)

	var foundAuthor models.Author
	if err := s.db.Where("username = ?", username).First(&foundAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Author %s not found", username)
			return models.Author{}, http.StatusNotFound, fmt.Errorf("author with username [%s] not found", username)
		}
		log.Printf("Error finding author: %v", err)
		return models.Author{}, http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundAuthor.Password), []byte(password)); err != nil {
		log.Printf("Invalid password for author %s", username)
		return models.Author{}, http.StatusUnauthorized, fmt.Errorf("invalid password")
	}

	log.Printf("Author %s authorized successfully", username)
	return foundAuthor, http.StatusOK, nil
}

func (s *Storage) DeleteAuthor(authorId uint32) (int, error) {
	log.Printf("Deleting author with ID: %d", authorId)

	if err := s.db.Where("id = ?", authorId).Delete(&models.Author{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Author with ID %d not found", authorId)
			return http.StatusConflict, err
		}
		log.Printf("Error deleting author: %v", err)
		return http.StatusInternalServerError, err
	}
	if err := s.db.Where("author_id = ?", authorId).Delete(&models.Article{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No articles found for author with ID %d", authorId)
			return http.StatusOK, nil
		}
		log.Printf("Error deleting articles for author with ID %d: %v", authorId, err)
		return http.StatusInternalServerError, err
	}
	log.Printf("Author with ID %d and their articles deleted successfully", authorId)
	return http.StatusOK, nil
}

func (s *Storage) GetAllArticles() ([]models.Article, error) {
	log.Println("Retrieving all articles")

	var foundArticles []models.Article
	if err := s.db.Preload("Author").Find(&foundArticles).Error; err != nil {
		log.Printf("Error retrieving articles: %v", err)
		return nil, err
	}
	log.Printf("Found %d articles", len(foundArticles))
	return foundArticles, nil
}

func (s *Storage) CreateArticle(article models.Article) error {
	log.Printf("Creating article: %s", article.Title)

	if err := s.db.Create(&article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		return err
	}
	log.Println("Article created successfully")
	return nil
}

func (s *Storage) DeleteArticle(articleId uint32) (int, error) {
	log.Printf("Deleting article with ID: %d", articleId)

	if err := s.db.Where("id = ?", articleId).Delete(&models.Article{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Article with ID %d not found", articleId)
			return http.StatusNotFound, err
		}
		log.Printf("Error deleting article: %v", err)
		return http.StatusInternalServerError, err
	}
	log.Printf("Article with ID %d deleted successfully", articleId)
	return http.StatusOK, nil
}
