package squarespace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/birddigital/store.adrienbird.net/internal/config"
	"github.com/birddigital/store.adrienbird.net/pkg/models"
)

type Client struct {
	baseURL     string
	siteID      string
	apiKey      string
	accessToken string
	httpClient  *http.Client
}

func NewClient(cfg *config.SquarespaceConfig) *Client {
	return &Client{
		baseURL:     cfg.BaseURL,
		siteID:      cfg.SiteID,
		apiKey:      cfg.APIKey,
		accessToken: cfg.AccessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "store.adrienbird.net/1.0")

	// Add authentication
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	} else if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	return c.httpClient.Do(req)
}

func (c *Client) decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiError models.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %s - %s", apiError.Type, apiError.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// Products API

func (c *Client) GetProducts(options ...ProductOption) ([]models.Product, *models.Pagination, error) {
	opts := &ProductOptions{}
	for _, opt := range options {
		opt(opts)
	}

	endpoint := "/1.0/commerce/products"
	if opts.SiteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/products", opts.SiteID)
	}

	// Add query parameters
	if opts.Limit > 0 || opts.Offset > 0 || opts.Category != "" || opts.Tag != "" {
		endpoint += "?"
		params := []string{}
		if opts.Limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
		}
		if opts.Offset > 0 {
			params = append(params, fmt.Sprintf("offset=%d", opts.Offset))
		}
		if opts.Category != "" {
			params = append(params, fmt.Sprintf("category=%s", opts.Category))
		}
		if opts.Tag != "" {
			params = append(params, fmt.Sprintf("tag=%s", opts.Tag))
		}
		for i, param := range params {
			if i > 0 {
				endpoint += "&"
			}
			endpoint += param
		}
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var response struct {
		Result    []models.Product `json:"result"`
		Pagination *models.Pagination `json:"pagination,omitempty"`
	}

	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, nil, err
	}

	return response.Result, response.Pagination, nil
}

func (c *Client) GetProduct(productID string) (*models.Product, error) {
	endpoint := fmt.Sprintf("/1.0/commerce/products/%s", productID)
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/products/%s", c.siteID, productID)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := c.decodeResponse(resp, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *Client) GetProductVariants(productID string) ([]models.ProductVariant, error) {
	product, err := c.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	return product.Products, nil
}

// Orders API

func (c *Client) GetOrders(options ...OrderOption) ([]models.Order, *models.Pagination, error) {
	opts := &OrderOptions{}
	for _, opt := range options {
		opt(opts)
	}

	endpoint := "/1.0/commerce/orders"
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/orders", c.siteID)
	}

	// Add query parameters
	if opts.Limit > 0 || opts.Offset > 0 || opts.Status != "" || opts.CustomerID != "" {
		endpoint += "?"
		params := []string{}
		if opts.Limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
		}
		if opts.Offset > 0 {
			params = append(params, fmt.Sprintf("offset=%d", opts.Offset))
		}
		if opts.Status != "" {
			params = append(params, fmt.Sprintf("status=%s", opts.Status))
		}
		if opts.CustomerID != "" {
			params = append(params, fmt.Sprintf("customerId=%s", opts.CustomerID))
		}
		for i, param := range params {
			if i > 0 {
				endpoint += "&"
			}
			endpoint += param
		}
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var response struct {
		Result    []models.Order `json:"result"`
		Pagination *models.Pagination `json:"pagination,omitempty"`
	}

	if err := c.decodeResponse(resp, &response); err != nil {
		return nil, nil, err
	}

	return response.Result, response.Pagination, nil
}

func (c *Client) GetOrder(orderID string) (*models.Order, error) {
	endpoint := fmt.Sprintf("/1.0/commerce/orders/%s", orderID)
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/orders/%s", c.siteID, orderID)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := c.decodeResponse(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *Client) CreateOrder(order *models.Order) (*models.Order, error) {
	endpoint := "/1.0/commerce/orders"
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/orders", c.siteID)
	}

	resp, err := c.makeRequest("POST", endpoint, order)
	if err != nil {
		return nil, err
	}

	var createdOrder models.Order
	if err := c.decodeResponse(resp, &createdOrder); err != nil {
		return nil, err
	}

	return &createdOrder, nil
}

// Inventory API

func (c *Client) GetInventory(productID string) (*models.ProductStock, error) {
	endpoint := fmt.Sprintf("/1.0/commerce/inventory/%s", productID)
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/inventory/%s", c.siteID, productID)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var inventory models.ProductStock
	if err := c.decodeResponse(resp, &inventory); err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (c *Client) UpdateInventory(productID string, quantity int) error {
	endpoint := fmt.Sprintf("/1.0/commerce/inventory/%s", productID)
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/inventory/%s", c.siteID, productID)
	}

	payload := map[string]interface{}{
		"quantity": quantity,
	}

	resp, err := c.makeRequest("PATCH", endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to update inventory with status %d", resp.StatusCode)
	}

	return nil
}

// Profiles API

func (c *Client) GetCustomerProfile(customerID string) (*models.Address, error) {
	endpoint := fmt.Sprintf("/1.0/commerce/profiles/%s", customerID)
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/profiles/%s", c.siteID, customerID)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var profile models.Address
	if err := c.decodeResponse(resp, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Health check

func (c *Client) HealthCheck() error {
	endpoint := "/1.0/commerce/products"
	if c.siteID != "" {
		endpoint = fmt.Sprintf("/1.0/commerce/sites/%s/products", c.siteID)
	}

	// Just try to fetch one product to check API connectivity
	endpoint += "?limit=1"

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API health check failed with status %d", resp.StatusCode)
	}

	return nil
}