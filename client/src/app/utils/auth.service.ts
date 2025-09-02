const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export class AuthService {
  // Store token in localStorage
  static setToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  // Get token from localStorage
  static getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('token');
    }
    return null;
  }

  // Remove token from localStorage
  static removeToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  // Check if user is authenticated
  static isAuthenticated(): boolean {
    const token = this.getToken();
    return !!token;
  }

  // Get auth headers
  static getAuthHeaders(): HeadersInit {
    const token = this.getToken();
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    };
  }

  // Login user
  static async login(credentials: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    const data = await response.json();

    if (data.success && data.data?.token) {
      this.setToken(data.data.token);
    }

    return data;
  }

  // Register user
  static async register(userData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    });

    const data = await response.json();
    return data;
  }

  // Logout user
  static logout(): void {
    this.removeToken();
  }

  // Get user profile
  static async getProfile(): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/api/auth/profile`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return await response.json();
  }

  // Update user profile
  static async updateProfile(userData: any): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/api/auth/profile`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(userData),
    });

    return await response.json();
  }
}