export interface SquarespaceConfig {
  apiUrl: string;
  apiKey?: string;
  accessToken?: string;
  siteId: string;
  environment: 'development' | 'production';
}

export const squarespaceConfig: SquarespaceConfig = {
  // Squarespace API base URL
  apiUrl: 'https://api.squarespace.com',

  // Your site identifier (get from Squarespace admin)
  siteId: process.env.SQUARESPACE_SITE_ID || '',

  // API credentials (set these in environment variables)
  apiKey: process.env.SQUARESPACE_API_KEY,
  accessToken: process.env.SQUARESPACE_ACCESS_TOKEN,

  // Environment
  environment: process.env.NODE_ENV === 'production' ? 'production' : 'development'
};

// API endpoints
export const SQUARESPACE_ENDPOINTS = {
  // Commerce APIs
  PRODUCTS: '/1.0/commerce/products',
  ORDERS: '/1.0/commerce/orders',
  INVENTORY: '/1.0/commerce/inventory',
  TRANSACTIONS: '/1.0/commerce/transactions',
  PROFILES: '/1.0/commerce/profiles',

  // Webhooks
  WEBHOOKS: '/1.0/webhooks/subscriptions'
} as const;