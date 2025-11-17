import { squarespaceConfig, SQUARESPACE_ENDPOINTS } from '../../config/squarespace.config';
import {
  Product,
  Order,
  SquarespaceApiResponse,
  SquarespaceError
} from '../types/squarespace';

class SquarespaceClient {
  private baseUrl: string;
  private apiKey?: string;
  private accessToken?: string;
  private siteId: string;

  constructor() {
    this.baseUrl = squarespaceConfig.apiUrl;
    this.apiKey = squarespaceConfig.apiKey;
    this.accessToken = squarespaceConfig.accessToken;
    this.siteId = squarespaceConfig.siteId;

    if (!this.siteId) {
      throw new Error('SQUARESPACE_SITE_ID environment variable is required');
    }
  }

  private async makeRequest<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'User-Agent': 'store.adrienbird.net/1.0',
      ...options.headers,
    };

    // Add authentication
    if (this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    } else if (this.apiKey) {
      headers['Authorization'] = `Bearer ${this.apiKey}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const errorData: SquarespaceError = await response.json();
        throw new SquarespaceApiError(errorData.message, response.status, errorData);
      }

      return await response.json();
    } catch (error) {
      if (error instanceof SquarespaceApiError) {
        throw error;
      }
      throw new SquarespaceApiError(
        `Network error: ${error instanceof Error ? error.message : 'Unknown error'}`,
        0
      );
    }
  }

  // Products API
  async getProducts(options?: {
    limit?: number;
    offset?: number;
    category?: string;
    tag?: string;
  }): Promise<SquarespaceApiResponse<Product>> {
    const params = new URLSearchParams();

    if (options?.limit) params.append('limit', options.limit.toString());
    if (options?.offset) params.append('offset', options.offset.toString());
    if (options?.category) params.append('category', options.category);
    if (options?.tag) params.append('tag', options.tag);

    const query = params.toString();
    const endpoint = `${SQUARESPACE_ENDPOINTS.PRODUCTS}${query ? `?${query}` : ''}`;

    return this.makeRequest<SquarespaceApiResponse<Product>>(endpoint);
  }

  async getProduct(productId: string): Promise<Product> {
    return this.makeRequest<Product>(`${SQUARESPACE_ENDPOINTS.PRODUCTS}/${productId}`);
  }

  async getProductVariants(productId: string): Promise<Product['products']> {
    const product = await this.getProduct(productId);
    return product.products;
  }

  // Orders API
  async getOrders(options?: {
    limit?: number;
    offset?: number;
    status?: string;
    customerId?: string;
  }): Promise<SquarespaceApiResponse<Order>> {
    const params = new URLSearchParams();

    if (options?.limit) params.append('limit', options.limit.toString());
    if (options?.offset) params.append('offset', options.offset.toString());
    if (options?.status) params.append('status', options.status);
    if (options?.customerId) params.append('customerId', options.customerId);

    const query = params.toString();
    const endpoint = `${SQUARESPACE_ENDPOINTS.ORDERS}${query ? `?${query}` : ''}`;

    return this.makeRequest<SquarespaceApiResponse<Order>>(endpoint);
  }

  async getOrder(orderId: string): Promise<Order> {
    return this.makeRequest<Order>(`${SQUARESPACE_ENDPOINTS.ORDERS}/${orderId}`);
  }

  async createOrder(orderData: Partial<Order>): Promise<Order> {
    return this.makeRequest<Order>(SQUARESPACE_ENDPOINTS.ORDERS, {
      method: 'POST',
      body: JSON.stringify(orderData),
    });
  }

  // Inventory API
  async getInventory(productId: string): Promise<any> {
    return this.makeRequest<any>(`${SQUARESPACE_ENDPOINTS.INVENTORY}/${productId}`);
  }

  async updateInventory(productId: string, quantity: number): Promise<any> {
    return this.makeRequest<any>(`${SQUARESPACE_ENDPOINTS.INVENTORY}/${productId}`, {
      method: 'PATCH',
      body: JSON.stringify({
        quantity,
      }),
    });
  }

  // Profiles API
  async getCustomerProfile(customerId: string): Promise<any> {
    return this.makeRequest<any>(`${SQUARESPACE_ENDPOINTS.PROFILES}/${customerId}`);
  }

  // Health check
  async healthCheck(): Promise<boolean> {
    try {
      await this.getProducts({ limit: 1 });
      return true;
    } catch (error) {
      return false;
    }
  }
}

export class SquarespaceApiError extends Error {
  public readonly statusCode: number;
  public readonly details?: any;

  constructor(message: string, statusCode: number, details?: any) {
    super(message);
    this.name = 'SquarespaceApiError';
    this.statusCode = statusCode;
    this.details = details;
  }
}

// Singleton instance
export const squarespaceClient = new SquarespaceClient();
export default squarespaceClient;