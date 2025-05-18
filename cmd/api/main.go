package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	_ "todo/docs"
	"todo/internal/config"
	"todo/internal/models"
	"todo/internal/server"
)

// @title Todo API
// @version 1.0
// @description A RESTful API for managing todos
// @host localhost:3035
// @BasePath /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create server
	srv := server.NewServer(db)

	// Start server
	port := strconv.Itoa(cfg.ServerPort)
	log.Printf("Server starting on port %s", port)
	err = srv.Start(":" + port)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
