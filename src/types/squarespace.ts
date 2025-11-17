export interface Product {
  id: string;
  type: 'PRODUCT';
  variantId: string;
  customForm?: {
    formId: string;
    fields: CustomFormField[];
  };
  categories?: string[];
  tags?: string[];
  products: ProductVariant[];
  relatedProducts?: RelatedProduct[];
  seoData?: SeoData;
  systemData: SystemData;
}

export interface ProductVariant {
  id: string;
  sku: string;
  name: string;
  description?: string;
  images: ProductImage[];
  pricing: {
    basePrice?: Money;
    compareAtPrice?: Money;
    salePrice?: Money;
    onSale: boolean;
  };
  stock: {
    trackInventory: boolean;
    quantity?: number;
    allowBackorder: boolean;
    unlimited: boolean;
  };
  visibility: 'PUBLIC' | 'PRIVATE' | 'HIDDEN';
  attributes?: ProductAttribute[];
  variants?: VariantOption[];
}

export interface ProductImage {
  assetId: string;
  url: string;
  description?: string;
  mimeType: string;
  width: number;
  height: number;
}

export interface Money {
  value: string;
  currency: string;
}

export interface CustomFormField {
  fieldId: string;
  type: string;
  label: string;
  required: boolean;
    choices?: string[];
  validation?: {
    minLength?: number;
    maxLength?: number;
    pattern?: string;
  };
}

export interface ProductAttribute {
  name: string;
  value: string;
}

export interface VariantOption {
  name: string;
  option: string;
}

export interface RelatedProduct {
  productId: string;
  variantId: string;
}

export interface SeoData {
  title?: string;
  description?: string;
  slug?: string;
  image?: string;
}

export interface SystemData {
  createdOn: number;
  modifiedOn: number;
  publishedOn: number;
}

// Order Types
export interface Order {
  id: string;
  orderNumber: string;
  customerId?: string;
  email: string;
  billingAddress: Address;
  shippingAddress?: Address;
  lineItems: OrderLineItem[];
  totals: OrderTotals;
  status: OrderStatus;
  fulfillments: OrderFulfillment[];
  systemData: SystemData;
}

export interface OrderLineItem {
  productId: string;
  variantId: string;
  sku: string;
  productName: string;
  variantName?: string;
  quantity: number;
  unitPrice: Money;
  totalPrice: Money;
  customizations?: OrderCustomization[];
}

export interface OrderCustomization {
  fieldName: string;
  value: string;
}

export interface OrderTotals {
  subtotal: Money;
  tax: Money;
  shipping: Money;
  discount: Money;
  total: Money;
}

export type OrderStatus =
  | 'PENDING'
  | 'CONFIRMED'
  | 'PROCESSING'
  | 'PARTIALLY_FULFILLED'
  | 'FULFILLED'
  | 'CANCELED'
  | 'REFUNDED';

export interface OrderFulfillment {
  id: string;
  type: 'SHIPPING' | 'PICKUP' | 'DIGITAL';
  status: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
  trackingInfo?: TrackingInfo;
  lineItems: string[]; // line item IDs
}

export interface TrackingInfo {
  carrier: string;
  trackingNumber: string;
  trackingUrl?: string;
}

export interface Address {
  firstName: string;
  lastName: string;
  company?: string;
  addressLine1: string;
  addressLine2?: string;
  city: string;
  state?: string;
  postalCode: string;
  country: string;
  phone?: string;
}

// API Response Types
export interface SquarespaceApiResponse<T> {
  result: T[];
  pagination?: {
    nextPage?: string;
    prevPage?: string;
    totalResults?: number;
  };
}

export interface SquarespaceError {
  type: string;
  message: string;
  details?: any;
}