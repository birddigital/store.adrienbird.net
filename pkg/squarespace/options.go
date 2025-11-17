package squarespace

type ProductOptions struct {
	SiteID   string
	Limit    int
	Offset   int
	Category string
	Tag      string
}

type ProductOption func(*ProductOptions)

func WithProductSiteID(siteID string) ProductOption {
	return func(opts *ProductOptions) {
		opts.SiteID = siteID
	}
}

func WithProductLimit(limit int) ProductOption {
	return func(opts *ProductOptions) {
		opts.Limit = limit
	}
}

func WithProductOffset(offset int) ProductOption {
	return func(opts *ProductOptions) {
		opts.Offset = offset
	}
}

func WithProductCategory(category string) ProductOption {
	return func(opts *ProductOptions) {
		opts.Category = category
	}
}

func WithProductTag(tag string) ProductOption {
	return func(opts *ProductOptions) {
		opts.Tag = tag
	}
}

type OrderOptions struct {
	SiteID     string
	Limit      int
	Offset     int
	Status     string
	CustomerID string
}

type OrderOption func(*OrderOptions)

func WithOrderSiteID(siteID string) OrderOption {
	return func(opts *OrderOptions) {
		opts.SiteID = siteID
	}
}

func WithOrderLimit(limit int) OrderOption {
	return func(opts *OrderOptions) {
		opts.Limit = limit
	}
}

func WithOrderOffset(offset int) OrderOption {
	return func(opts *OrderOptions) {
		opts.Offset = offset
	}
}

func WithOrderStatus(status string) OrderOption {
	return func(opts *OrderOptions) {
		opts.Status = status
	}
}

func WithOrderCustomerID(customerID string) OrderOption {
	return func(opts *OrderOptions) {
		opts.CustomerID = customerID
	}
}