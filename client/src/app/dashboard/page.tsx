'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import ProtectedRoute from '../components/ProtectedRoute';
import { useAuth } from '../contexts/auth.context';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <ProtectedRoute>
      {/* Header/Navbar */}
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <h1 className="text-xl font-bold text-gray-900">Blog Platform</h1>
              </div>
              <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                <Link href="/dashboard" className="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                  Dashboard
                </Link>
                <Link href="/posts" className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                  Posts
                </Link>
                <Link href="/posts/create" className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                  Create Post
                </Link>
                <Link href="/profile" className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                  Profile
                </Link>
              </div>
            </div>
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <button
                  onClick={handleLogout}
                  className="relative inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Logout
                </button>
              </div>
              <div className="ml-3 relative">
                <div className="flex items-center">
                  <div className="text-sm text-right mr-3">
                    <div className="font-medium text-gray-900">{user?.first_name} {user?.last_name}</div>
                    <div className="text-gray-500">@{user?.username}</div>
                  </div>
                  <div className="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
                    <span className="text-indigo-800 font-medium">
                      {user?.first_name?.charAt(0)}{user?.last_name?.charAt(0)}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0">
            <div className="border-4 border-dashed border-gray-200 rounded-lg h-96 flex items-center justify-center">
              <div className="text-center">
                <h2 className="text-2xl font-bold text-gray-900 mb-4">Welcome to your Dashboard</h2>
                <p className="text-gray-600 mb-6">This is your personalized dashboard where you can manage your blog posts and settings.</p>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <Link href="/posts" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
                    <h3 className="text-lg font-medium text-gray-900">Your Posts</h3>
                    <p className="mt-2 text-gray-600">Manage your blog posts</p>
                  </Link>
                  <Link href="/posts/create" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
                    <h3 className="text-lg font-medium text-gray-900">Create New Post</h3>
                    <p className="mt-2 text-gray-600">Write a new blog post</p>
                  </Link>
                  <Link href="/profile" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
                    <h3 className="text-lg font-medium text-gray-900">Settings</h3>
                    <p className="mt-2 text-gray-600">Update your profile and preferences</p>
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}