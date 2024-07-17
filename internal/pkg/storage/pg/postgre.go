package pg

import (
	"articleModule/internal/pkg/models"
	"articleModule/internal/utils/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
