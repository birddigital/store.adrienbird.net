package handlers

import (
	"net/http"
	"time"

	"github.com/birddigital/store.adrienbird.net/internal/config"
	"github.com/birddigital/store.adrienbird.net/pkg/models"
	"github.com/birddigital/store.adrienbird.net/pkg/squarespace"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	cfg    *config.Config
	client *squarespace.Client
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		cfg:    cfg,
		client: squarespace.NewClient(&cfg.Squarespace),
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	healthResponse := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    make(map[string]models.Health),
	}

	// Check Squarespace API connectivity
	squarespaceCheck := models.Health{
		Status: "healthy",
	}
	start := time.Now()

	if err := h.client.HealthCheck(); err != nil {
		squarespaceCheck.Status = "unhealthy"
		squarespaceCheck.Message = err.Error()
	} else {
		squarespaceCheck.Latency = time.Since(start)
	}
	healthResponse.Checks["squarespace_api"] = squarespaceCheck

	// Check configuration
	configCheck := models.Health{
		Status: "healthy",
	}
	if h.cfg.Squarespace.SiteID == "" {
		configCheck.Status = "warning"
		configCheck.Message = "SQUARESPACE_SITE_ID not configured"
	}
	if h.cfg.Squarespace.APIKey == "" && h.cfg.Squarespace.AccessToken == "" {
		configCheck.Status = "error"
		configCheck.Message = "No Squarespace authentication configured"
	}
	healthResponse.Checks["configuration"] = configCheck

	// Determine overall status
	if configCheck.Status == "error" {
		healthResponse.Status = "unhealthy"
	} else if squarespaceCheck.Status == "unhealthy" || configCheck.Status == "warning" {
		healthResponse.Status = "degraded"
	}

	// Set HTTP status based on health
	statusCode := http.StatusOK
	if healthResponse.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if healthResponse.Status == "degraded" {
		statusCode = http.StatusOK // Still OK but with warnings
	}

	c.JSON(statusCode, healthResponse)
}