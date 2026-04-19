package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/api/routes"
	"github.com/mythicalsystems/nemesis-backend/internal/config"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("[ERR][LOGGER] Failed to initialize logger: %s", err)
	}
	defer logger.GetLogger().Sync()

	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[ERR][CONFIG] Failed to load config: %s", err)
	}
	log.Println("[OK][CONFIG] Config loaded successfully")

	// Initialize Database
	if err := database.Init(cfg); err != nil {
		log.Fatalf("[ERR][DATABASE] Failed to initialize database: %s", err)
	}
	log.Println("[OK][DATABASE] Database initialized successfully")

	// Setup Router
	r := gin.Default()

	// Register Routes
	routes.Register(r)

	// Start Server
	log.Println("[OK][SERVER] Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[ERR][SERVER] Failed to start server: %s", err)
	}
}
