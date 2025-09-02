'use client';

import { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import { ApiResponse, AuthContextType, LoginRequest, LoginResponse, RegisterRequest, User } from '../types/auth.types';
import { AuthService } from '../utils/auth.service';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // Check if user is authenticated on initial load
  useEffect(() => {
    const initializeAuth = async () => {
      setIsLoading(true);
      const storedToken = AuthService.getToken();

      if (storedToken) {
        setToken(storedToken);

        try {
          const profileResponse = await AuthService.getProfile();
          if (profileResponse.success) {
            setUser(profileResponse.data);
          } else {
            // Token is invalid, remove it
            AuthService.removeToken();
            setToken(null);
          }
        } catch (error) {
          // Error fetching profile, remove token
          AuthService.removeToken();
          setToken(null);
        }
      }

      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    try {
      const response = await AuthService.login(credentials);

      if (response.success && response.data?.token) {
        setToken(response.data.token);
        setUser(response.data.user);
        AuthService.setToken(response.data.token);
      }

      return response;
    } catch (error) {
      return {
        success: false,
        error: 'An error occurred during login',
      };
    }
  };

  const register = async (userData: RegisterRequest): Promise<ApiResponse<LoginResponse>> => {
    try {
      const response = await AuthService.register(userData);
      return response;
    } catch (error) {
      return {
        success: false,
        error: 'An error occurred during registration',
      };
    }
  };

  const logout = (): void => {
    AuthService.logout();
    setToken(null);
    setUser(null);
  };

  const value = {
    user,
    token,
    login,
    register,
    logout,
    isAuthenticated: !!user && !!token,
    isLoading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

// Custom hook to use the auth context
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}