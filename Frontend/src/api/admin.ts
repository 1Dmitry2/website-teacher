import { ADMIN_TOKEN_KEY } from '../utils/adminAuth';
import { API_BASE_URL, ApiClient } from './client';
import type { AuthResponse, LoginRequest } from './client';

export interface AdminProfile {
  id: number;
  email: string;
  created_at?: string;
}

export interface DashboardSummary {
  message: string;
  generated_at?: string;
}

export interface ResetPasswordPayload {
  token: string;
  newPassword: string;
}

export interface Block {
  id: string;
  page: string;
  pages: string[];
  type: 'text' | 'slider' | 'gallery' | 'video' | 'text-with-image' | 'document';
  content: Record<string, any>;
  display_order: number;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: string;
  title: string;
  content: string;
  images: string[];
  videos?: string[];
  pages: string[];
  is_published: boolean;
  alignment?: 'left' | 'center' | 'right' | 'full-width';
  title_position?: 'top' | 'bottom' | 'left' | 'right';
  content_position?: 'top' | 'bottom' | 'left' | 'right';
  created_at: string;
  updated_at: string;
}

export interface GalleryItem {
  id: string;
  image_url: string;
  title?: string;
  description?: string;
  text?: string;
  pages: string[];
  created_at: string;
  updated_at: string;
}

export interface SliderItem {
  id: string;
  image_url: string;
  title?: string;
  description?: string;
  pages: string[];
  display_order: number;
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

export interface User {
  id: number;
  email: string;
  created_at?: string;
  is_admin?: boolean;
}

class AdminApiClient extends ApiClient {
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/admin/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.data!;
  }

  async forgotPassword(email: string): Promise<void> {
    await this.request('/admin/forgot-password', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  async resetPassword(payload: ResetPasswordPayload): Promise<void> {
    await this.request('/admin/reset-password', {
      method: 'POST',
      body: JSON.stringify({
        token: payload.token,
        new_password: payload.newPassword,
      }),
    });
  }

  async me(): Promise<AdminProfile> {
    const response = await this.request<AdminProfile>('/admin/me', {
      method: 'GET',
    });
    return response.data!;
  }

  async dashboard(): Promise<DashboardSummary> {
    const response = await this.request<DashboardSummary>('/admin/dashboard', {
      method: 'GET',
    });
    return response.data!;
  }

  async getBlocks(): Promise<Block[]> {
    const response = await this.request<Block[]>('/admin/blocks', { method: 'GET' });
    return response.data || [];
  }

  async createBlock(block: Omit<Block, 'id' | 'created_at' | 'updated_at'>): Promise<Block> {
    const response = await this.request<Block>('/admin/blocks', {
      method: 'POST',
      body: JSON.stringify(block),
    });
    return response.data!;
  }

  async updateBlock(id: string, block: Partial<Block>): Promise<Block> {
    const response = await this.request<Block>(`/admin/blocks/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(block),
    });
    return response.data!;
  }

  async deleteBlock(id: string): Promise<void> {
    await this.request(`/admin/blocks/${id}`, { method: 'DELETE' });
  }

  async reorderBlocks(items: { id: string; display_order: number }[]): Promise<void> {
    await this.request('/admin/blocks/reorder', {
      method: 'PATCH',
      body: JSON.stringify({ items }),
    });
  }

  async getPosts(): Promise<Post[]> {
    const response = await this.request<Post[]>('/admin/posts', { method: 'GET' });
    return response.data || [];
  }

  async getPost(id: string): Promise<Post> {
    const response = await this.request<Post>(`/admin/posts/${id}`, { method: 'GET' });
    return response.data!;
  }

  async createPost(post: Omit<Post, 'id' | 'created_at' | 'updated_at'>): Promise<Post> {
    const response = await this.request<Post>('/admin/posts', {
      method: 'POST',
      body: JSON.stringify(post),
    });
    return response.data!;
  }

  async updatePost(id: string, post: Partial<Post>): Promise<Post> {
    const response = await this.request<Post>(`/admin/posts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(post),
    });
    return response.data!;
  }

  async deletePost(id: string): Promise<void> {
    await this.request(`/admin/posts/${id}`, { method: 'DELETE' });
  }

  async getGallery(): Promise<GalleryItem[]> {
    const response = await this.request<GalleryItem[]>('/admin/gallery', { method: 'GET' });
    return response.data || [];
  }

  async createGalleryItem(item: Omit<GalleryItem, 'id' | 'created_at' | 'updated_at'>): Promise<GalleryItem> {
    const response = await this.request<GalleryItem>('/admin/gallery', {
      method: 'POST',
      body: JSON.stringify(item),
    });
    return response.data!;
  }

  async updateGalleryItem(id: string, item: Partial<GalleryItem>): Promise<GalleryItem> {
    const response = await this.request<GalleryItem>(`/admin/gallery/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(item),
    });
    return response.data!;
  }

  async deleteGalleryItem(id: string): Promise<void> {
    await this.request(`/admin/gallery/${id}`, { method: 'DELETE' });
  }

  async getSlider(): Promise<SliderItem[]> {
    const response = await this.request<SliderItem[]>('/admin/slider', { method: 'GET' });
    return response.data || [];
  }

  async createSliderItem(item: Omit<SliderItem, 'id' | 'created_at' | 'updated_at'>): Promise<SliderItem> {
    const response = await this.request<SliderItem>('/admin/slider', {
      method: 'POST',
      body: JSON.stringify(item),
    });
    return response.data!;
  }

  async updateSliderItem(id: string, item: Partial<SliderItem>): Promise<SliderItem> {
    const response = await this.request<SliderItem>(`/admin/slider/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(item),
    });
    return response.data!;
  }

  async deleteSliderItem(id: string): Promise<void> {
    await this.request(`/admin/slider/${id}`, { method: 'DELETE' });
  }

  async reorderSlider(items: { id: string; display_order: number }[]): Promise<void> {
    await this.request('/admin/slider/reorder', {
      method: 'PATCH',
      body: JSON.stringify({ items }),
    });
  }

  async getComments(): Promise<Comment[]> {
    const response = await this.request<Comment[]>('/admin/comments', { method: 'GET' });
    return response.data || [];
  }

  async replyToComment(commentId: string, text: string): Promise<Comment> {
    const response = await this.request<Comment>(`/admin/comments/${commentId}/reply`, {
      method: 'POST',
      body: JSON.stringify({ text }),
    });
    return response.data!;
  }

  async deleteComment(id: string): Promise<void> {
    await this.request(`/admin/comments/${id}`, { method: 'DELETE' });
  }

  async getUsers(): Promise<User[]> {
    const response = await this.request<User[]>('/admin/users', { method: 'GET' });
    return response.data || [];
  }

  async getUserById(userId: number): Promise<User & { comments_count?: number }> {
    const response = await this.request<User & { comments_count?: number }>(`/admin/users/${userId}`, {
      method: 'GET',
    });
    return response.data!;
  }

  async banUser(userId: number, banned: boolean): Promise<void> {
    await this.request(`/admin/users/${userId}/ban`, {
      method: 'PATCH',
      body: JSON.stringify({ banned }),
    });
  }

  async uploadFile(file: File): Promise<string> {
    const token = this.getToken();
    if (!token) {
      throw new Error('Токен авторизации не найден. Пожалуйста, войдите в систему.');
    }

    const formData = new FormData();
    formData.append('file', file);

    const headers: Record<string, string> = {
      'Authorization': `Bearer ${token}`,
    };

    const url = `${this.baseURL}/admin/upload`;
    console.log('Uploading to:', url);
    console.log('Token present:', !!token);

    const response = await fetch(url, {
      method: 'POST',
      headers: headers as HeadersInit,
      body: formData,
    });

    console.log('Upload response status:', response.status);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Upload failed' }));
      console.error('Upload error:', errorData);
      throw new Error(errorData.error || errorData.data?.error || `HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.data?.url || data.url || '';
  }
}

export const adminApi = new AdminApiClient(API_BASE_URL, ADMIN_TOKEN_KEY);

