// API Client for Xephyr Backend
// Handles requests to the Go backend API

import { ApiResponse, ApiError } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Default request timeout in milliseconds
const DEFAULT_TIMEOUT = 30000;

// Custom error class for API errors
export class ApiClientError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, string[]>,
    public statusCode?: number
  ) {
    super(message);
    this.name = 'ApiClientError';
  }
}

// Request configuration interface
interface RequestConfig extends RequestInit {
  timeout?: number;
}

// Build full URL
function buildUrl(path: string): string {
  // If path already starts with http, use it as-is
  if (path.startsWith('http')) {
    return path;
  }
  
  // Remove leading slash from path if present
  const cleanPath = path.startsWith('/') ? path.slice(1) : path;
  
  // Remove trailing slash from base URL if present
  const cleanBase = API_BASE_URL.endsWith('/') 
    ? API_BASE_URL.slice(0, -1) 
    : API_BASE_URL;
  
  return `${cleanBase}/api/v1/${cleanPath}`;
}

// Get default headers
function getDefaultHeaders(): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };

  // Add auth token if available (can be extended with real auth)
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // Add organization ID if available
  const orgId = typeof window !== 'undefined' ? localStorage.getItem('organization_id') : null;
  if (orgId) {
    headers['X-Organization-Id'] = orgId;
  }

  return headers;
}

// Handle API response
async function handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
  // Check if response is JSON
  const contentType = response.headers.get('content-type');
  const isJson = contentType?.includes('application/json');

  let data: any;
  try {
    data = isJson ? await response.json() : await response.text();
  } catch (e) {
    data = null;
  }

  if (!response.ok) {
    // Handle different error types
    if (data?.error) {
      throw new ApiClientError(
        data.error.code || 'UNKNOWN_ERROR',
        data.error.message || 'An error occurred',
        data.error.details,
        response.status
      );
    }
    
    // HTTP status-based errors
    const message = data?.message || response.statusText || 'Request failed';
    throw new ApiClientError(
      `HTTP_${response.status}`,
      message,
      undefined,
      response.status
    );
  }

  // Return the API response format
  if (data?.success !== undefined) {
    return data as ApiResponse<T>;
  }

  // Wrap raw data in standard format
  return {
    success: true,
    data: data as T,
    meta: {
      timestamp: new Date().toISOString(),
      requestId: `req_${Date.now()}`,
    },
  };
}

// Make a request with timeout
async function makeRequest<T>(
  url: string,
  config: RequestConfig = {}
): Promise<ApiResponse<T>> {
  const { timeout = DEFAULT_TIMEOUT, ...fetchConfig } = config;
  
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...fetchConfig,
      signal: controller.signal,
    });
    clearTimeout(timeoutId);
    return handleResponse<T>(response);
  } catch (error) {
    clearTimeout(timeoutId);
    
    if (error instanceof ApiClientError) {
      throw error;
    }
    
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new ApiClientError('TIMEOUT', 'Request timed out');
      }
      throw new ApiClientError('NETWORK_ERROR', error.message);
    }
    
    throw new ApiClientError('UNKNOWN_ERROR', 'An unknown error occurred');
  }
}

// ==================== HTTP METHODS ====================

/**
 * Make a GET request
 */
export async function get<T>(
  path: string,
  params?: Record<string, string | number | boolean | undefined>,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  // Build query string
  let url = buildUrl(path);
  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        searchParams.append(key, String(value));
      }
    });
    const queryString = searchParams.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
  }

  return makeRequest<T>(url, {
    ...config,
    method: 'GET',
    headers: {
      ...getDefaultHeaders(),
      ...config?.headers,
    },
  });
}

/**
 * Make a POST request
 */
export async function post<T>(
  path: string,
  body?: unknown,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  return makeRequest<T>(buildUrl(path), {
    ...config,
    method: 'POST',
    headers: {
      ...getDefaultHeaders(),
      ...config?.headers,
    },
    body: body ? JSON.stringify(body) : undefined,
  });
}

/**
 * Make a PATCH request
 */
export async function patch<T>(
  path: string,
  body?: unknown,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  return makeRequest<T>(buildUrl(path), {
    ...config,
    method: 'PATCH',
    headers: {
      ...getDefaultHeaders(),
      ...config?.headers,
    },
    body: body ? JSON.stringify(body) : undefined,
  });
}

/**
 * Make a PUT request
 */
export async function put<T>(
  path: string,
  body?: unknown,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  return makeRequest<T>(buildUrl(path), {
    ...config,
    method: 'PUT',
    headers: {
      ...getDefaultHeaders(),
      ...config?.headers,
    },
    body: body ? JSON.stringify(body) : undefined,
  });
}

/**
 * Make a DELETE request
 */
export async function del<T>(
  path: string,
  config?: RequestConfig
): Promise<ApiResponse<T>> {
  return makeRequest<T>(buildUrl(path), {
    ...config,
    method: 'DELETE',
    headers: {
      ...getDefaultHeaders(),
      ...config?.headers,
    },
  });
}

// ==================== API CLIENT OBJECT ====================

export const apiClient = {
  get,
  post,
  patch,
  put,
  delete: del,
};

export default apiClient;
