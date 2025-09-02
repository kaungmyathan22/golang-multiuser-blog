// User type matching the backend UserResponse model
export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  username: string;
  bio: string;
  avatar: string;
  is_active: boolean;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

// Login request type matching the backend UserLoginRequest model
export interface LoginRequest {
  email_or_username: string;
  password: string;
}

// Register request type matching the backend UserCreateRequest model
export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  username: string;
  password: string;
  bio?: string;
  avatar?: string;
}

// Login response type matching the backend AuthResponse model
export interface LoginResponse {
  user: User;
  token: string;
  token_type: string;
  expires_in: number;
  refresh_token?: string;
}

// API response wrapper matching the backend APIResponse model
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string | object;
}

// Auth context type
export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (credentials: LoginRequest) => Promise<ApiResponse<LoginResponse>>;
  register: (userData: RegisterRequest) => Promise<ApiResponse<LoginResponse>>;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}