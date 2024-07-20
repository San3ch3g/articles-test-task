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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	return dsn
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MustNewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Author{})
	if err != nil {
		log.Println("table 'authors' is already exist")
	}
	err = db.AutoMigrate(&models.Article{})
	if err != nil {
		log.Println("table 'articles' is already exist")
	}
	return db
}

func (s *Storage) RegisterAuthor(username, password string) (models.Author, int, error) {

	var foundAuthor models.Author
	if err := s.db.Where("username = ?", username).First(&foundAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newAuthor models.Author
			newAuthor.Username = username
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return models.Author{}, http.StatusInternalServerError, err
			}
			newAuthor.Password = string(hashedPassword)
			if err := s.db.Create(&newAuthor).Error; err != nil {
				return models.Author{}, http.StatusInternalServerError, err
			}
			return newAuthor, http.StatusCreated, nil
		} else {
			return models.Author{}, http.StatusConflict, err
		}
	}
	return models.Author{}, http.StatusConflict, fmt.Errorf("author with username [%s] already exists", username)
}

func (s *Storage) AuthorizeAuthor(username, password string) (models.Author, int, error) {
	var foundAuthor models.Author
	if err := s.db.Where("username = ?", username).First(&foundAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Author{}, http.StatusNotFound, fmt.Errorf("author with username [%s] not found", username)
		}
		return models.Author{}, http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundAuthor.Password), []byte(password)); err != nil {
		return models.Author{}, http.StatusUnauthorized, fmt.Errorf("invalid password")
	}

	return foundAuthor, http.StatusOK, nil
}

func (s *Storage) DeleteAuthor(authorId uint32) (int, error) {
	if err := s.db.Where("id = ?", authorId).Delete(&models.Author{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusConflict, err
		}
		return http.StatusInternalServerError, err
	}
	if err := s.db.Where("author_id = ?", authorId).Delete(&models.Article{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *Storage) GetAllArticles() ([]models.Article, error) {
	var foundArticles []models.Article
	if err := s.db.Preload("Author").Find(&foundArticles).Error; err != nil {
		return nil, err
	}
	return foundArticles, nil
}

func (s *Storage) CreateArticle(article models.Article) error {
	if err := s.db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteArticle(articleId uint32) (int, error) {
	if err := s.db.Where("id = ?", articleId).Delete(&models.Article{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
