'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useAuth } from '../../contexts/auth.context';
import { Post } from '../../types/blog.types';
import { BlogService } from '../../utils/blog.service';

export default function PostDetailPage({ params }: { params: { id: string } }) {
  const { isAuthenticated, user } = useAuth();
  const router = useRouter();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchPost();
  }, [params.id]);

  const fetchPost = async () => {
    setLoading(true);
    setError('');

    try {
      const response = await BlogService.getPostById(parseInt(params.id));

      if (response.success) {
        setPost(response.data);
      } else {
        setError(response.error || 'Failed to fetch post');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = () => {
    router.push(`/posts/${params.id}/edit`);
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this post?')) {
      try {
        const response = await BlogService.deletePost(parseInt(params.id));

        if (response.success) {
          router.push('/posts');
        } else {
          setError(response.error || 'Failed to delete post');
        }
      } catch (err) {
        setError('An unexpected error occurred');
      }
    }
  };

  const handlePublish = async () => {
    try {
      const response = await BlogService.publishPost(parseInt(params.id));

      if (response.success) {
        // Refresh the post data
        fetchPost();
      } else {
        setError(response.error || 'Failed to publish post');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    }
  };

  const handleUnpublish = async () => {
    try {
      const response = await BlogService.unpublishPost(parseInt(params.id));

      if (response.success) {
        // Refresh the post data
        fetchPost();
      } else {
        setError(response.error || 'Failed to unpublish post');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="bg-red-50 border-l-4 border-red-400 p-4">
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
      </div>
    );
  }

  if (!post) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900">Post not found</h1>
          <p className="mt-2 text-gray-600">The post you're looking for doesn't exist.</p>
          <Link
            href="/posts"
            className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Back to Posts
          </Link>
        </div>
      </div>
    );
  }

  const isAuthor = user?.id === post.author.id;
  const canEdit = isAuthenticated && (isAuthor || user?.is_admin);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Post Header */}
          <div className="bg-white shadow sm:rounded-lg mb-6">
            <div className="px-4 py-5 sm:px-6">
              <div className="flex justify-between items-start">
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">{post.title}</h1>
                  <div className="mt-2 flex items-center text-sm text-gray-500">
                    <span>By {post.author.first_name} {post.author.last_name}</span>
                    <span className="mx-2">•</span>
                    <time dateTime={post.created_at}>
                      {new Date(post.created_at).toLocaleDateString()}
                    </time>
                    {post.published_at && (
                      <>
                        <span className="mx-2">•</span>
                        <span>Published on {new Date(post.published_at).toLocaleDateString()}</span>
                      </>
                    )}
                  </div>
                </div>

                {canEdit && (
                  <div className="flex space-x-2">
                    {!post.is_published ? (
                      <button
                        onClick={handlePublish}
                        className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                      >
                        Publish
                      </button>
                    ) : (
                      <button
                        onClick={handleUnpublish}
                        className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-yellow-600 hover:bg-yellow-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
                      >
                        Unpublish
                      </button>
                    )}
                    <button
                      onClick={handleEdit}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    >
                      Edit
                    </button>
                    <button
                      onClick={handleDelete}
                      className="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                    >
                      Delete
                    </button>
                  </div>
                )}
              </div>

              {post.excerpt && (
                <p className="mt-3 text-lg text-gray-600">{post.excerpt}</p>
              )}

              <div className="mt-4 flex flex-wrap gap-2">
                {post.tags && post.tags.map((tag) => (
                  <span
                    key={tag.id}
                    className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    style={{ backgroundColor: `${tag.color}20`, color: tag.color }}
                  >
                    {tag.name}
                  </span>
                ))}
                {!post.is_published && (
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    Draft
                  </span>
                )}
              </div>
            </div>
          </div>

          {/* Post Content */}
          <div className="bg-white shadow sm:rounded-lg">
            <div className="px-4 py-5 sm:px-6">
              {post.featured_image && (
                <div className="mb-6">
                  <img
                    src={post.featured_image}
                    alt={post.title}
                    className="w-full h-64 object-cover rounded-lg"
                  />
                </div>
              )}

              <div className="prose max-w-none">
                <div dangerouslySetInnerHTML={{ __html: post.content.replace(/\n/g, '<br />') }} />
              </div>
            </div>
          </div>

          {/* Back to Posts Button */}
          <div className="mt-6">
            <Link
              href="/posts"
              className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              ← Back to Posts
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}