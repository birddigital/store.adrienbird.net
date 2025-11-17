package handlers

import (
	"net/http"
	"strconv"

	"github.com/birddigital/store.adrienbird.net/internal/config"
	"github.com/birddigital/store.adrienbird.net/pkg/models"
	"github.com/birddigital/store.adrienbird.net/pkg/squarespace"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	cfg    *config.Config
	client *squarespace.Client
}

func NewProductHandler(cfg *config.Config) *ProductHandler {
	return &ProductHandler{
		cfg:    cfg,
		client: squarespace.NewClient(&cfg.Squarespace),
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	category := c.Query("category")
	tag := c.Query("tag")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "invalid_parameter",
				Message: "Invalid limit parameter",
			},
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "invalid_parameter",
				Message: "Invalid offset parameter",
			},
		})
		return
	}

	// Build options
	options := []squarespace.ProductOption{
		squarespace.WithProductLimit(limit),
		squarespace.WithProductOffset(offset),
	}
	if category != "" {
		options = append(options, squarespace.WithProductCategory(category))
	}
	if tag != "" {
		options = append(options, squarespace.WithProductTag(tag))
	}

	// Fetch products from Squarespace
	products, pagination, err := h.client.GetProducts(options...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{
				Type:    "api_error",
				Message: "Failed to fetch products: " + err.Error(),
			},
		})
		return
	}

	// Build response
	response := models.APIResponse{
		Data:       products,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "missing_parameter",
				Message: "Product ID is required",
			},
		})
		return
	}

	// Fetch product from Squarespace
	product, err := h.client.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{
				Type:    "not_found",
				Message: "Product not found: " + err.Error(),
			},
		})
		return
	}

	response := models.APIResponse{
		Data: product,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductVariants(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "missing_parameter",
				Message: "Product ID is required",
			},
		})
		return
	}

	// Fetch product variants from Squarespace
	variants, err := h.client.GetProductVariants(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{
				Type:    "not_found",
				Message: "Product variants not found: " + err.Error(),
			},
		})
		return
	}

	response := models.APIResponse{
		Data: variants,
	}

	c.JSON(http.StatusOK, response)
}