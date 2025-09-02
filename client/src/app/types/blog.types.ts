// User type (imported from auth.types.ts)
import { User } from './auth.types';

// Tag type matching the backend Tag model
export interface Tag {
  id: number;
  name: string;
  slug: string;
  description: string;
  color: string;
  created_at: string;
  updated_at: string;
}

// Comment type matching the backend Comment model
export interface Comment {
  id: number;
  content: string;
  author: User;
  post_id: number;
  parent_id: number | null;
  is_approved: boolean;
  created_at: string;
  updated_at: string;
}

// Post type matching the backend Post model
export interface Post {
  id: number;
  title: string;
  content: string;
  slug: string;
  excerpt: string;
  featured_image: string;
  is_published: boolean;
  view_count: number;
  author: User;
  tags: Tag[];
  comments: Comment[];
  created_at: string;
  updated_at: string;
  published_at: string | null;
}

// Create post request type
export interface CreatePostRequest {
  title: string;
  content: string;
  excerpt: string;
  featured_image?: string;
  tag_ids?: number[];
}

// Update post request type
export interface UpdatePostRequest {
  title?: string;
  content?: string;
  excerpt?: string;
  featured_image?: string;
  tag_ids?: number[];
  is_published?: boolean;
}

// API response wrapper for posts
export interface PostApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string | object;
}

// Pagination metadata
export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

// Paginated response for posts
export interface PaginatedPostsResponse {
  success: boolean;
  message?: string;
  data: Post[];
  pagination: PaginationMeta;
  error?: string | object;
}