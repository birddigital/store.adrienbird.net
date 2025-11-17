package main

import (
	"fmt"
	"log"
	"os"

	"github.com/birddigital/store.adrienbird.net/internal/config"
	"github.com/birddigital/store.adrienbird.net/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Gin router
	router := gin.New()

	// Setup middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Setup CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://adrienbird.net")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	productHandler := handlers.NewProductHandler(cfg)
	orderHandler := handlers.NewOrderHandler(cfg)
	healthHandler := handlers.NewHealthHandler(cfg)

	// Setup routes
	api := router.Group("/api/v1")
	{
		// Product routes
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProduct)
		api.GET("/products/:id/variants", productHandler.GetProductVariants)

		// Order routes
		api.GET("/orders", orderHandler.GetOrders)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.POST("/orders", orderHandler.CreateOrder)

		// Health check
		api.GET("/health", healthHandler.Health)
	}

	// Add root health endpoint
	router.GET("/health", healthHandler.Health)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Store.AdrienBird.net API",
			"version": "1.0.0",
			"status":  "healthy",
		})
	})

	// Start server
	log.Printf("Starting server on port %d", cfg.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}