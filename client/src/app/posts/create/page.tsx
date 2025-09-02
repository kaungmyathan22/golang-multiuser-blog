'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import ProtectedRoute from '../../components/ProtectedRoute';
import { useAuth } from '../../contexts/auth.context';
import { CreatePostRequest } from '../../types/blog.types';
import { BlogService } from '../../utils/blog.service';

export default function CreatePostPage() {
  const { user } = useAuth();
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [excerpt, setExcerpt] = useState('');
  const [featuredImage, setFeaturedImage] = useState('');
  const [tags, setTags] = useState('');
  const [isPublished, setIsPublished] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      // Convert tags string to array of numbers (assuming tag IDs)
      const tagIds = tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag !== '')
        .map(tag => parseInt(tag))
        .filter(tag => !isNaN(tag));

      // Determine status based on isPublished flag
      const status: 'published' | 'draft' | 'archived' = isPublished ? 'published' : 'draft';

      const postData: CreatePostRequest = {
        title,
        content,
        excerpt,
        featured_image: featuredImage || undefined,
        tag_ids: tagIds.length > 0 ? tagIds : undefined,
        status,
      };

      const response = await BlogService.createPost(postData);

      if (response.success) {
        setSuccess('Post created successfully!');
        // Reset form
        setTitle('');
        setContent('');
        setExcerpt('');
        setFeaturedImage('');
        setTags('');
        setIsPublished(false);

        // Redirect to posts page after a short delay
        setTimeout(() => {
          router.push('/posts');
        }, 1500);
      } else {
        setError(response.error || 'Failed to create post');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0">
            <div className="bg-white shadow sm:rounded-lg">
              <div className="px-4 py-5 sm:px-6">
                <h1 className="text-2xl font-bold text-gray-900">Create New Post</h1>
                <p className="mt-1 text-sm text-gray-500">
                  Fill in the details below to create a new blog post.
                </p>
              </div>

              {error && (
                <div className="bg-red-50 border-l-4 border-red-400 p-4 mx-4 sm:mx-6">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                      </svg>
                    </div>
                    <div className="ml-3">
                      <p className="text-sm text-red-700">
                        {error}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              {success && (
                <div className="bg-green-50 border-l-4 border-green-400 p-4 mx-4 sm:mx-6">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <svg className="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                      </svg>
                    </div>
                    <div className="ml-3">
                      <p className="text-sm text-green-700">
                        {success}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              <div className="border-t border-gray-200">
                <form onSubmit={handleSubmit} className="px-4 py-5 sm:p-6">
                  <div className="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
                    <div className="sm:col-span-6">
                      <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                        Title
                      </label>
                      <div className="mt-1">
                        <input
                          type="text"
                          name="title"
                          id="title"
                          value={title}
                          onChange={(e) => setTitle(e.target.value)}
                          required
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                    </div>

                    <div className="sm:col-span-6">
                      <label htmlFor="excerpt" className="block text-sm font-medium text-gray-700">
                        Excerpt
                      </label>
                      <div className="mt-1">
                        <textarea
                          id="excerpt"
                          name="excerpt"
                          rows={3}
                          value={excerpt}
                          onChange={(e) => setExcerpt(e.target.value)}
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                    </div>

                    <div className="sm:col-span-6">
                      <label htmlFor="content" className="block text-sm font-medium text-gray-700">
                        Content
                      </label>
                      <div className="mt-1">
                        <textarea
                          id="content"
                          name="content"
                          rows={10}
                          value={content}
                          onChange={(e) => setContent(e.target.value)}
                          required
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                    </div>

                    <div className="sm:col-span-6">
                      <label htmlFor="featured-image" className="block text-sm font-medium text-gray-700">
                        Featured Image URL
                      </label>
                      <div className="mt-1">
                        <input
                          type="text"
                          name="featured-image"
                          id="featured-image"
                          value={featuredImage}
                          onChange={(e) => setFeaturedImage(e.target.value)}
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                    </div>

                    <div className="sm:col-span-6">
                      <label htmlFor="tags" className="block text-sm font-medium text-gray-700">
                        Tag IDs (comma separated)
                      </label>
                      <div className="mt-1">
                        <input
                          type="text"
                          name="tags"
                          id="tags"
                          value={tags}
                          onChange={(e) => setTags(e.target.value)}
                          placeholder="1, 2, 3"
                          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                        />
                      </div>
                    </div>

                    <div className="sm:col-span-6">
                      <div className="flex items-center">
                        <input
                          id="published"
                          name="published"
                          type="checkbox"
                          checked={isPublished}
                          onChange={(e) => setIsPublished(e.target.checked)}
                          className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        />
                        <label htmlFor="published" className="ml-2 block text-sm text-gray-900">
                          Publish immediately
                        </label>
                      </div>
                    </div>
                  </div>

                  <div className="mt-6 flex items-center justify-end space-x-3">
                    <button
                      type="button"
                      onClick={() => router.push('/posts')}
                      className="inline-flex justify-center py-2 px-4 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={loading}
                      className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                    >
                      {loading ? (
                        <span className="flex items-center">
                          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          Creating...
                        </span>
                      ) : (
                        'Create Post'
                      )}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}