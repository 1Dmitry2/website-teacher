const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
}

export interface User {
  id: number;
  email: string;
}

export interface ApiResponse<T> {
  message: string;
  data?: T;
  error?: string;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const token = localStorage.getItem('token');
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (options.headers) {
      Object.assign(headers, options.headers);
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, {
        ...options,
        headers: headers as HeadersInit,
      });

      const data: ApiResponse<T> = await response.json();

      if (!response.ok) {
        const errorMessage = data.error || 
          (data.data && typeof data.data === 'object' && 'error' in data.data 
            ? (data.data as any).error 
            : `HTTP error! status: ${response.status}`);
        throw new Error(errorMessage);
      }

      return data;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Network error occurred');
    }
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/login-users', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data!;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/sessions', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data!;
  }

  async getProfile(): Promise<User> {
    const response = await this.request<User>('/profile', {
      method: 'GET',
    });
    return response.data!;
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
