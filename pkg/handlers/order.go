package handlers

import (
	"net/http"
	"strconv"

	"github.com/birddigital/store.adrienbird.net/internal/config"
	"github.com/birddigital/store.adrienbird.net/pkg/models"
	"github.com/birddigital/store.adrienbird.net/pkg/squarespace"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	cfg    *config.Config
	client *squarespace.Client
}

func NewOrderHandler(cfg *config.Config) *OrderHandler {
	return &OrderHandler{
		cfg:    cfg,
		client: squarespace.NewClient(&cfg.Squarespace),
	}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	status := c.Query("status")
	customerID := c.Query("customerId")

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
	options := []squarespace.OrderOption{
		squarespace.WithOrderLimit(limit),
		squarespace.WithOrderOffset(offset),
	}
	if status != "" {
		options = append(options, squarespace.WithOrderStatus(status))
	}
	if customerID != "" {
		options = append(options, squarespace.WithOrderCustomerID(customerID))
	}

	// Fetch orders from Squarespace
	orders, pagination, err := h.client.GetOrders(options...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{
				Type:    "api_error",
				Message: "Failed to fetch orders: " + err.Error(),
			},
		})
		return
	}

	// Build response
	response := models.APIResponse{
		Data:       orders,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "missing_parameter",
				Message: "Order ID is required",
			},
		})
		return
	}

	// Fetch order from Squarespace
	order, err := h.client.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{
				Type:    "not_found",
				Message: "Order not found: " + err.Error(),
			},
		})
		return
	}

	response := models.APIResponse{
		Data: order,
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "invalid_request",
				Message: "Invalid order data: " + err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if order.Email == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "validation_error",
				Message: "Email is required",
			},
		})
		return
	}

	if len(order.LineItems) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{
				Type:    "validation_error",
				Message: "At least one line item is required",
			},
		})
		return
	}

	// Create order in Squarespace
	createdOrder, err := h.client.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{
				Type:    "creation_error",
				Message: "Failed to create order: " + err.Error(),
			},
		})
		return
	}

	response := models.APIResponse{
		Data: createdOrder,
	}

	c.JSON(http.StatusCreated, response)
}