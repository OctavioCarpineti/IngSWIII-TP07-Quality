import React, { useEffect, useState } from 'react';
import { postService } from '../../services/postService';
import { Post } from '../../types';
import './PostList.css';

interface PostListProps {
    currentUserId: number;
    onRefresh?: boolean;
    onViewPost?: (postId: number) => void;
}

export const PostList: React.FC<PostListProps> = ({ currentUserId, onRefresh, onViewPost }) => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const loadPosts = async () => {
        try {
            setLoading(true);
            const data = await postService.getAllPosts();
            setPosts(data);
            setError('');
        } catch (err: any) {
            setError('Error al cargar posts');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadPosts();
    }, [onRefresh]);

    const handleDelete = async (postId: number) => {
        if (!window.confirm('¿Estás seguro de eliminar este post?')) {
            return;
        }

        try {
            await postService.deletePost(postId, currentUserId);
            loadPosts();
        } catch (err: any) {
            alert(err.response?.data?.error || 'Error al eliminar post');
        }
    };

    if (loading) {
        return <div className="loading">Cargando posts...</div>;
    }

    if (error) {
        return <div className="error">{error}</div>;
    }

    if (posts.length === 0) {
        return <div className="no-posts">No hay posts todavía. ¡Crea el primero!</div>;
    }

    return (
        <div className="post-list">
            <h2>Posts</h2>
            {posts.map((post) => (
                <div
                    key={post.id}
                    className="post-card"
                    onClick={() => onViewPost && onViewPost(post.id)}
                    style={{ cursor: onViewPost ? 'pointer' : 'default' }}
                >
                    <div className="post-header">
                        <h3>{post.title}</h3>
                        <span className="post-author">por @{post.username}</span>
                    </div>
                    <p className="post-content">{post.content}</p>
                    <div className="post-footer">
            <span className="post-date">
              {new Date(post.created_at).toLocaleDateString()}
            </span>
                        {post.user_id === currentUserId && (
                            <button
                                className="delete-btn"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleDelete(post.id);
                                }}
                            >
                                Eliminar
                            </button>
                        )}
                    </div>
                </div>
            ))}
        </div>
    );
};
