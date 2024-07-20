package main

import (
	_ "articleModule/docs"
	"articleModule/internal/di"
	"articleModule/internal/pkg/handlers"
	"articleModule/internal/utils/config"
	"fmt"
)

// @title Article Swagger
// @version 1.0
// @description This is an API for test task
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.NewConfig()
	cfg.InitENV()

	container := di.New(cfg)
	db := container.GetDB()

	postgresDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database connection: %v", err))
	}
	if err := postgresDB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	storage := container.GetSQLStorage()
	server := handlers.NewServer(storage, cfg)
	server.InitSwagger()
	err = server.Run(cfg.ServerPort)
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
