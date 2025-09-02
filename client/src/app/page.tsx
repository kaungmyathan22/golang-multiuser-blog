'use client';

import Link from 'next/link';
import { useAuth } from './contexts/auth.context';

export default function Home() {
  const { user, isAuthenticated } = useAuth();

  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Multi-User Blog</h1>
          <p className="text-lg text-gray-600 mb-8">A modern blogging platform for multiple authors</p>

          {isAuthenticated ? (
            <div className="mb-8">
              <p className="text-xl mb-4">Welcome back, {user?.first_name}!</p>
              <div className="flex gap-4 justify-center">
                <Link
                  href="/dashboard"
                  className="rounded-full border border-solid border-transparent transition-colors flex items-center justify-center bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc] font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5"
                >
                  Go to Dashboard
                </Link>
                <Link
                  href="/profile"
                  className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5"
                >
                  View Profile
                </Link>
              </div>
            </div>
          ) : (
            <div className="mb-8">
              <p className="text-xl mb-4">Get started with our blogging platform</p>
              <div className="flex gap-4 justify-center">
                <Link
                  href="/login"
                  className="rounded-full border border-solid border-transparent transition-colors flex items-center justify-center bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc] font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5"
                >
                  Sign In
                </Link>
                <Link
                  href="/register"
                  className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5"
                >
                  Register
                </Link>
              </div>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full max-w-4xl">
          <div className="border border-gray-200 rounded-lg p-6 text-center">
            <h3 className="text-xl font-semibold mb-2">Write Posts</h3>
            <p className="text-gray-600">Create and publish your blog posts with our easy-to-use editor.</p>
          </div>
          <div className="border border-gray-200 rounded-lg p-6 text-center">
            <h3 className="text-xl font-semibold mb-2">Engage</h3>
            <p className="text-gray-600">Connect with readers through comments and discussions.</p>
          </div>
          <div className="border border-gray-200 rounded-lg p-6 text-center">
            <h3 className="text-xl font-semibold mb-2">Manage</h3>
            <p className="text-gray-600">Track your posts, comments, and reader engagement.</p>
          </div>
        </div>
      </main>

      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        <Link
          href="/dashboard"
          className="flex items-center gap-2 hover:underline hover:underline-offset-4"
        >
          Dashboard
        </Link>
        <Link
          href="/profile"
          className="flex items-center gap-2 hover:underline hover:underline-offset-4"
        >
          Profile
        </Link>
        {isAuthenticated ? (
          <Link
            href="/login"
            className="flex items-center gap-2 hover:underline hover:underline-offset-4"
          >
            Sign Out
          </Link>
        ) : (
          <>
            <Link
              href="/login"
              className="flex items-center gap-2 hover:underline hover:underline-offset-4"
            >
              Sign In
            </Link>
            <Link
              href="/register"
              className="flex items-center gap-2 hover:underline hover:underline-offset-4"
            >
              Register
            </Link>
          </>
        )}
      </footer>
    </div>
  );
}