import { TOKEN_KEY } from '../utils/auth';

export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

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

export interface Post {
  id: string;
  title: string;
  content: string;
  images: string[];
  pages: string[];
  is_published: boolean;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  post_id: string;
  user_id: number;
  reply_to?: string;
  text: string;
  is_admin: boolean;
  created_at: string;
}

export interface ApiResponse<T> {
  message: string;
  data?: T;
  error?: string;
}

export class ApiClient {
  protected baseURL: string;
  private tokenKey: string | null;

  constructor(baseURL: string, tokenKey: string | null = TOKEN_KEY) {
    this.baseURL = baseURL;
    this.tokenKey = tokenKey;
  }

  protected getToken(): string | null {
    if (!this.tokenKey) {
      return null;
    }
    return localStorage.getItem(this.tokenKey);
  }

  protected async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const token = this.getToken();
    
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

      const rawText = await response.text();
      let data: ApiResponse<T>;

      try {
        data = rawText ? JSON.parse(rawText) : { message: 'ok' };
      } catch {
        throw new Error(`Некорректный ответ сервера: ${rawText.slice(0, 200) || 'пусто'}`);
      }

      if (!response.ok) {
        const errorMessage =
          data.error ||
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

  // Public CMS endpoints
  async getPageBlocks(page: string): Promise<any[]> {
    const response = await this.request<any[]>(`/pages/${page}`, { method: 'GET' });
    return response.data || [];
  }

  async getPosts(): Promise<Post[]> {
    const response = await this.request<Post[]>('/posts', { method: 'GET' });
    return response.data || [];
  }

  async getPost(id: string): Promise<any> {
    const response = await this.request<any>(`/posts/${id}`, { method: 'GET' });
    return response.data!;
  }

  async getPostComments(postId: string): Promise<Comment[]> {
    const response = await this.request<Comment[]>(`/posts/${postId}/comments`, { method: 'GET' });
    return response.data || [];
  }

  async createComment(postId: string, text: string): Promise<Comment> {
    const response = await this.request<Comment>(`/posts/${postId}/comments`, {
      method: 'POST',
      body: JSON.stringify({ text }),
    });
    return response.data!;
  }

  async replyToComment(commentId: string, text: string): Promise<Comment> {
    const response = await this.request<Comment>(`/comments/${commentId}/reply`, {
      method: 'POST',
      body: JSON.stringify({ text }),
    });
    return response.data!;
  }

  async deleteComment(commentId: string): Promise<void> {
    await this.request(`/comments/${commentId}`, { method: 'DELETE' });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
