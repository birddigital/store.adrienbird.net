package models

import "time"

type Product struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	VariantID    string        `json:"variantId"`
	CustomForm   *CustomForm   `json:"customForm,omitempty"`
	Categories   []string      `json:"categories,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Products     []ProductVariant `json:"products"`
	RelatedProducts []RelatedProduct `json:"relatedProducts,omitempty"`
	SeoData      *SeoData      `json:"seoData,omitempty"`
	SystemData   SystemData    `json:"systemData"`
}

type ProductVariant struct {
	ID          string            `json:"id"`
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Images      []ProductImage    `json:"images"`
	Pricing     ProductPricing    `json:"pricing"`
	Stock       ProductStock      `json:"stock"`
	Visibility  string            `json:"visibility"`
	Attributes  []ProductAttribute `json:"attributes,omitempty"`
	Variants    []VariantOption   `json:"variants,omitempty"`
}

type ProductImage struct {
	AssetID     string `json:"assetId"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type ProductPricing struct {
	BasePrice      *Money `json:"basePrice,omitempty"`
	CompareAtPrice *Money `json:"compareAtPrice,omitempty"`
	SalePrice      *Money `json:"salePrice,omitempty"`
	OnSale         bool   `json:"onSale"`
}

type ProductStock struct {
	TrackInventory bool   `json:"trackInventory"`
	Quantity       *int   `json:"quantity,omitempty"`
	AllowBackorder bool   `json:"allowBackorder"`
	Unlimited      bool   `json:"unlimited"`
}

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type VariantOption struct {
	Name   string `json:"name"`
	Option string `json:"option"`
}

type RelatedProduct struct {
	ProductID  string `json:"productId"`
	VariantID  string `json:"variantId"`
}

type SeoData struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Image       string `json:"image,omitempty"`
}

type SystemData struct {
	CreatedOn   int64 `json:"createdOn"`
	ModifiedOn  int64 `json:"modifiedOn"`
	PublishedOn int64 `json:"publishedOn"`
}

type CustomForm struct {
	FormID  string           `json:"formId"`
	Fields  []CustomFormField `json:"fields"`
}

type CustomFormField struct {
	FieldID    string     `json:"fieldId"`
	Type       string     `json:"type"`
	Label      string     `json:"label"`
	Required   bool       `json:"required"`
	Choices    []string   `json:"choices,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

type Validation struct {
	MinLength *int    `json:"minLength,omitempty"`
	MaxLength *int    `json:"maxLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`
}

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Order struct {
	ID            string          `json:"id"`
	OrderNumber   string          `json:"orderNumber"`
	CustomerID    *string         `json:"customerId,omitempty"`
	Email         string          `json:"email"`
	BillingAddress Address         `json:"billingAddress"`
	ShippingAddress *Address      `json:"shippingAddress,omitempty"`
	LineItems     []OrderLineItem `json:"lineItems"`
	Totals        OrderTotals     `json:"totals"`
	Status        string          `json:"status"`
	Fulfillments  []OrderFulfillment `json:"fulfillments"`
	SystemData    SystemData      `json:"systemData"`
}

type OrderLineItem struct {
	ProductID      string               `json:"productId"`
	VariantID      string               `json:"variantId"`
	SKU            string               `json:"sku"`
	ProductName    string               `json:"productName"`
	VariantName    *string              `json:"variantName,omitempty"`
	Quantity       int                  `json:"quantity"`
	UnitPrice      Money                `json:"unitPrice"`
	TotalPrice     Money                `json:"totalPrice"`
	Customizations []OrderCustomization `json:"customizations,omitempty"`
}

type OrderCustomization struct {
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

type OrderTotals struct {
	Subtotal Money `json:"subtotal"`
	Tax      Money `json:"tax"`
	Shipping Money `json:"shipping"`
	Discount Money `json:"discount"`
	Total    Money `json:"total"`
}

type OrderFulfillment struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Status       string         `json:"status"`
	TrackingInfo *TrackingInfo  `json:"trackingInfo,omitempty"`
	LineItems    []string       `json:"lineItems"`
}

type TrackingInfo struct {
	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"trackingNumber"`
	TrackingURL    *string `json:"trackingUrl,omitempty"`
}

type Address struct {
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Company      *string `json:"company,omitempty"`
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	City         string  `json:"city"`
	State        *string `json:"state,omitempty"`
	PostalCode   string  `json:"postalCode"`
	Country      string  `json:"country"`
	Phone        *string `json:"phone,omitempty"`
}

type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type APIError struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Pagination struct {
	NextPage     *string `json:"nextPage,omitempty"`
	PrevPage     *string `json:"prevPage,omitempty"`
	TotalResults *int    `json:"totalResults,omitempty"`
}

// Health check models
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Checks    map[string]Health `json:"checks,omitempty"`
}

type Health struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}