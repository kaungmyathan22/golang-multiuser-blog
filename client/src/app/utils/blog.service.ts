import {
    CreatePostRequest,
    PaginatedPostsResponse,
    UpdatePostRequest
} from '../types/blog.types';
import { AuthService } from './auth.service';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export class BlogService {
  // Get all posts (published)
  static async getPublishedPosts(page: number = 1, perPage: number = 10): Promise<PaginatedPostsResponse> {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/posts/published?page=${page}&per_page=${perPage}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to fetch published posts',
        data: [],
        pagination: {
          page: 1,
          per_page: 10,
          total: 0,
          total_pages: 0,
        },
      };
    }
  }

  // Get all posts (for authenticated users)
  static async getPosts(page: number = 1, perPage: number = 10): Promise<PaginatedPostsResponse> {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/posts?page=${page}&per_page=${perPage}`,
        {
          method: 'GET',
          headers: AuthService.getAuthHeaders(),
        }
      );
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to fetch posts',
        data: [],
        pagination: {
          page: 1,
          per_page: 10,
          total: 0,
          total_pages: 0,
        },
      };
    }
  }

  // Get post by ID
  static async getPostById(id: number): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/${id}`, {
        method: 'GET',
        headers: AuthService.getAuthHeaders(),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to fetch post',
      };
    }
  }

  // Get post by slug
  static async getPostBySlug(slug: string): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/slug/${slug}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to fetch post',
      };
    }
  }

  // Create post
  static async createPost(postData: CreatePostRequest): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts`, {
        method: 'POST',
        headers: AuthService.getAuthHeaders(),
        body: JSON.stringify(postData),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to create post',
      };
    }
  }

  // Update post
  static async updatePost(id: number, postData: UpdatePostRequest): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/${id}`, {
        method: 'PUT',
        headers: AuthService.getAuthHeaders(),
        body: JSON.stringify(postData),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to update post',
      };
    }
  }

  // Delete post
  static async deletePost(id: number): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/${id}`, {
        method: 'DELETE',
        headers: AuthService.getAuthHeaders(),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to delete post',
      };
    }
  }

  // Publish post
  static async publishPost(id: number): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/${id}/publish`, {
        method: 'POST',
        headers: AuthService.getAuthHeaders(),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to publish post',
      };
    }
  }

  // Unpublish post
  static async unpublishPost(id: number): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/posts/${id}/unpublish`, {
        method: 'POST',
        headers: AuthService.getAuthHeaders(),
      });
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to unpublish post',
      };
    }
  }

  // Search posts
  static async searchPosts(query: string, page: number = 1, perPage: number = 10): Promise<PaginatedPostsResponse> {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/posts/search?q=${encodeURIComponent(query)}&page=${page}&per_page=${perPage}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      return await response.json();
    } catch (error) {
      return {
        success: false,
        error: 'Failed to search posts',
        data: [],
        pagination: {
          page: 1,
          per_page: 10,
          total: 0,
          total_pages: 0,
        },
      };
    }
  }
}