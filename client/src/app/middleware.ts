import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

// Define which paths require authentication
const protectedPaths = ['/dashboard', '/profile'];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Check if the current path requires authentication
  const isProtectedPath = protectedPaths.some(path => pathname.startsWith(path));

  if (isProtectedPath) {
    // In a real implementation, you would check for a valid session/token here
    // For now, we'll just demonstrate the concept
    const token = request.cookies.get('token');

    if (!token) {
      // Redirect to login if no token
      return NextResponse.redirect(new URL('/login', request.url));
    }

    // In a real implementation, you would validate the token here
    // For now, we'll assume it's valid
  }

  return NextResponse.next();
}

// Configure which paths the middleware should run on
export const config = {
  matcher: ['/dashboard/:path*', '/profile/:path*'],
};